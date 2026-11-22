// agents-server owns the agents database on the Windows host. agents-push and
// agents-pull start it: locally on Windows and WSL, over SSH from other devices.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"dotfiles/src/helpers/agentstore"

	_ "modernc.org/sqlite"
)

const defaultDB = "D:/Store/agents.db"

const schema = `
CREATE TABLE IF NOT EXISTS files (
	id         INTEGER PRIMARY KEY,
	device     TEXT    NOT NULL,
	path       TEXT    NOT NULL,           -- relative to ~/.claude, with / separators
	kind       TEXT    NOT NULL,           -- 'lines' (.jsonl) or 'blob'
	project    TEXT,                       -- project folder name under projects/
	session    TEXT,                       -- session id the file belongs to
	size       INTEGER NOT NULL,
	hash       TEXT    NOT NULL,           -- sha256 of the whole file
	mtime      INTEGER NOT NULL,           -- modified time on the device, unix ms
	tail_open  INTEGER NOT NULL DEFAULT 0, -- 1 when the file does not end with a newline
	content    BLOB,                       -- file bytes, blobs only
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	UNIQUE (device, path)
);
CREATE INDEX IF NOT EXISTS files_session ON files (session);
CREATE INDEX IF NOT EXISTS files_project ON files (device, project);

CREATE TABLE IF NOT EXISTS lines (
	file_id INTEGER NOT NULL REFERENCES files (id) ON DELETE CASCADE,
	n       INTEGER NOT NULL,              -- line number, from 0
	json    TEXT    NOT NULL,
	PRIMARY KEY (file_id, n)
) WITHOUT ROWID;

CREATE VIEW IF NOT EXISTS messages AS
	SELECT f.device, f.project, f.session, f.path, l.n, l.json
	FROM lines l JOIN files f ON f.id = l.file_id;
`

const (
	batchFiles = 2000
	batchBytes = 256 << 20
)

var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

type storedFile struct {
	id       int64
	kind     string
	size     int64
	hash     string
	mtime    int64
	tailOpen bool
}

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "agents-server: "+format+"\n", args...)
}

func main() {
	defaultPath := defaultDB
	if env := os.Getenv("AGENTS_DB"); env != "" {
		defaultPath = env
	}
	listen := flag.Bool("listen", false, "accept one connection on a local port instead of using standard input and output")
	dbPath := flag.String("db", defaultPath, "database file")
	flag.Parse()

	db, dbErr := openDB(*dbPath)
	if db != nil {
		defer db.Close()
	}

	if !*listen {
		serve(db, dbErr, os.Stdin, os.Stdout, "")
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logf("failed to listen: %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		panic(err)
	}
	token := hex.EncodeToString(tokenBytes)
	fmt.Printf("%s %d %s\n", agentstore.ReadyPrefix, listener.Addr().(*net.TCPAddr).Port, token)

	_ = listener.(*net.TCPListener).SetDeadline(time.Now().Add(2 * time.Minute))
	conn, err := listener.Accept()
	if err != nil {
		logf("no connection: %v", err)
		os.Exit(1)
	}
	listener.Close()
	defer conn.Close()
	serve(db, dbErr, conn, conn, token)
}

func openDB(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(600000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to prepare %s: %w", path, err)
	}
	return db, nil
}

func serve(db *sql.DB, dbErr error, r io.Reader, w io.Writer, token string) {
	conn := agentstore.NewConn(r, w)

	var hello agentstore.Hello
	if err := conn.Receive(&hello); err != nil {
		logf("failed to read request: %v", err)
		return
	}

	fail := func(message string) {
		logf("%s", message)
		_ = conn.Send(agentstore.State{Err: "agents-server: " + message})
	}
	switch {
	case token != "" && hello.Token != token:
		fail("wrong token")
		return
	case hello.Version != agentstore.ProtocolVersion:
		fail(fmt.Sprintf("version %d does not match server version %d, recompile scripts", hello.Version, agentstore.ProtocolVersion))
		return
	case hello.Device == "":
		fail("no device given")
		return
	case hello.Op != "push" && hello.Op != "pull":
		fail("unknown operation " + hello.Op)
		return
	case dbErr != nil:
		fail(dbErr.Error())
		return
	}

	state, err := loadState(db, hello.Device)
	if err != nil {
		fail("failed to read stored files: " + err.Error())
		return
	}
	if err := conn.Send(state); err != nil {
		logf("failed to send state: %v", err)
		return
	}

	if hello.Op == "push" {
		handlePush(db, conn, hello.Device)
	} else {
		handlePull(db, conn, hello.Device)
	}
}

func loadState(db *sql.DB, device string) (agentstore.State, error) {
	rows, err := db.Query(`SELECT path, kind, size, hash, mtime FROM files WHERE device = ?`, device)
	if err != nil {
		return agentstore.State{}, err
	}
	defer rows.Close()

	var state agentstore.State
	for rows.Next() {
		var file agentstore.FileState
		if err := rows.Scan(&file.Path, &file.Kind, &file.Size, &file.Hash, &file.Mtime); err != nil {
			return agentstore.State{}, err
		}
		state.Files = append(state.Files, file)
	}
	return state, rows.Err()
}

func handlePush(db *sql.DB, conn *agentstore.Conn, device string) {
	var result agentstore.PushResult
	tx, err := db.Begin()
	if err != nil {
		_ = conn.Send(agentstore.PushResult{Err: "agents-server: " + err.Error()})
		return
	}
	pendingFiles, pendingBytes := 0, 0

	commit := func() error {
		if err := tx.Commit(); err != nil {
			return err
		}
		pendingFiles, pendingBytes = 0, 0
		tx, err = db.Begin()
		return err
	}

	for {
		var upload agentstore.Upload
		if err := conn.Receive(&upload); err != nil {
			// The device went away. Keep what arrived completely.
			logf("push from %s stopped: %v", device, err)
			_ = tx.Commit()
			return
		}
		if upload.Done {
			break
		}

		outcome, err := applyUploadSafely(tx, device, &upload)
		if err != nil {
			result.Errors = append(result.Errors, upload.Path+": "+err.Error())
			continue
		}
		switch outcome {
		case "touched":
			result.Touched++
		case "appended":
			result.Appended++
		case "replaced":
			result.Replaced++
		case "added":
			result.Added++
		}
		result.Bytes += int64(len(upload.Data))

		pendingFiles++
		pendingBytes += len(upload.Data)
		if pendingFiles >= batchFiles || pendingBytes >= batchBytes {
			if err := commit(); err != nil {
				_ = conn.Send(agentstore.PushResult{Err: "agents-server: failed to save: " + err.Error()})
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		_ = conn.Send(agentstore.PushResult{Err: "agents-server: failed to save: " + err.Error()})
		return
	}
	_ = conn.Send(result)
}

// applyUploadSafely applies one upload so that it is either fully stored or
// not stored at all.
func applyUploadSafely(tx *sql.Tx, device string, upload *agentstore.Upload) (string, error) {
	if _, err := tx.Exec(`SAVEPOINT upload`); err != nil {
		return "", err
	}
	outcome, err := applyUpload(tx, device, upload)
	if err != nil {
		_, _ = tx.Exec(`ROLLBACK TO upload`)
		_, _ = tx.Exec(`RELEASE upload`)
		return "", err
	}
	_, err = tx.Exec(`RELEASE upload`)
	return outcome, err
}

func applyUpload(tx *sql.Tx, device string, upload *agentstore.Upload) (string, error) {
	if !validPath(upload.Path) {
		return "", errors.New("invalid path")
	}

	existing, err := findFile(tx, device, upload.Path)
	if err != nil {
		return "", err
	}
	now := time.Now().UnixMilli()

	switch upload.Mode {
	case agentstore.ModeTouch:
		if existing == nil || existing.hash != upload.Hash {
			return "", errors.New("stored copy changed, push again")
		}
		_, err := tx.Exec(`UPDATE files SET mtime = ?, updated_at = ? WHERE id = ?`, upload.Mtime, now, existing.id)
		return "touched", err

	case agentstore.ModeAppend:
		if existing == nil || existing.kind != agentstore.KindLines || existing.size != upload.Base {
			return "", errors.New("stored copy changed, push again")
		}
		if existing.size+int64(len(upload.Data)) != upload.Size {
			return "", errors.New("size does not match")
		}
		next, err := lineCount(tx, existing.id)
		if err != nil {
			return "", err
		}
		tailOpen, err := insertLines(tx, existing.id, next, existing.tailOpen, upload.Data)
		if err != nil {
			return "", err
		}
		_, err = tx.Exec(`UPDATE files SET size = ?, hash = ?, mtime = ?, tail_open = ?, updated_at = ? WHERE id = ?`,
			upload.Size, upload.Hash, upload.Mtime, tailOpen, now, existing.id)
		return "appended", err

	case agentstore.ModeFull:
		if int64(len(upload.Data)) != upload.Size || hashOf(upload.Data) != upload.Hash {
			return "", errors.New("data does not match its hash")
		}
		project, session := idsFor(upload.Path)
		var content []byte
		if upload.Kind == agentstore.KindBlob {
			content = upload.Data
			if content == nil {
				content = []byte{}
			}
		}

		outcome := "added"
		var id int64
		if existing == nil {
			res, err := tx.Exec(`INSERT INTO files (device, path, kind, project, session, size, hash, mtime, tail_open, content, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?)`,
				device, upload.Path, upload.Kind, nullable(project), nullable(session), upload.Size, upload.Hash, upload.Mtime, content, now, now)
			if err != nil {
				return "", err
			}
			if id, err = res.LastInsertId(); err != nil {
				return "", err
			}
		} else {
			outcome = "replaced"
			id = existing.id
			if _, err := tx.Exec(`DELETE FROM lines WHERE file_id = ?`, id); err != nil {
				return "", err
			}
			if _, err := tx.Exec(`UPDATE files SET kind = ?, project = ?, session = ?, size = ?, hash = ?, mtime = ?, tail_open = 0, content = ?, updated_at = ? WHERE id = ?`,
				upload.Kind, nullable(project), nullable(session), upload.Size, upload.Hash, upload.Mtime, content, now, id); err != nil {
				return "", err
			}
		}

		if upload.Kind == agentstore.KindLines {
			tailOpen, err := insertLines(tx, id, 0, false, upload.Data)
			if err != nil {
				return "", err
			}
			if _, err := tx.Exec(`UPDATE files SET tail_open = ? WHERE id = ?`, tailOpen, id); err != nil {
				return "", err
			}
		}
		return outcome, nil
	}
	return "", errors.New("unknown mode " + upload.Mode)
}

func findFile(tx *sql.Tx, device string, path string) (*storedFile, error) {
	var file storedFile
	err := tx.QueryRow(`SELECT id, kind, size, hash, mtime, tail_open FROM files WHERE device = ? AND path = ?`, device, path).
		Scan(&file.id, &file.kind, &file.size, &file.hash, &file.mtime, &file.tailOpen)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func lineCount(tx *sql.Tx, fileID int64) (int64, error) {
	var count int64
	err := tx.QueryRow(`SELECT COALESCE(MAX(n) + 1, 0) FROM lines WHERE file_id = ?`, fileID).Scan(&count)
	return count, err
}

// insertLines stores data as lines starting at line number next. When the
// stored file ends without a newline, the first piece of data continues that
// last line. It returns whether the file still ends without a newline.
func insertLines(tx *sql.Tx, fileID int64, next int64, tailOpen bool, data []byte) (bool, error) {
	if len(data) == 0 {
		return tailOpen, nil
	}

	if tailOpen {
		piece := data
		index := bytes.IndexByte(data, '\n')
		if index >= 0 {
			piece = data[:index]
		}
		if _, err := tx.Exec(`UPDATE lines SET json = json || ? WHERE file_id = ? AND n = ?`, string(piece), fileID, next-1); err != nil {
			return false, err
		}
		if index < 0 {
			return true, nil
		}
		data = data[index+1:]
	}

	statement, err := tx.Prepare(`INSERT INTO lines (file_id, n, json) VALUES (?, ?, ?)`)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	for len(data) > 0 {
		index := bytes.IndexByte(data, '\n')
		if index < 0 {
			if _, err := statement.Exec(fileID, next, string(data)); err != nil {
				return false, err
			}
			return true, nil
		}
		if _, err := statement.Exec(fileID, next, string(data[:index])); err != nil {
			return false, err
		}
		data = data[index+1:]
		next++
	}
	return false, nil
}

func handlePull(db *sql.DB, conn *agentstore.Conn, device string) {
	var want agentstore.Want
	if err := conn.Receive(&want); err != nil {
		logf("failed to read pull request: %v", err)
		return
	}

	for _, path := range want.Paths {
		file, err := readFile(db, device, path)
		if err != nil {
			logf("skipped %s: %v", path, err)
			continue
		}
		if file == nil {
			continue
		}
		if err := conn.Send(file); err != nil {
			logf("pull to %s stopped: %v", device, err)
			return
		}
	}
	_ = conn.Send(agentstore.File{Done: true})
}

func readFile(db *sql.DB, device string, path string) (*agentstore.File, error) {
	var (
		id       int64
		size     int64
		hash     string
		tailOpen bool
		content  []byte
	)
	file := agentstore.File{Path: path}
	err := db.QueryRow(`SELECT id, kind, size, hash, mtime, tail_open, content FROM files WHERE device = ? AND path = ?`, device, path).
		Scan(&id, &file.Kind, &size, &hash, &file.Mtime, &tailOpen, &content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if file.Kind == agentstore.KindBlob {
		file.Data = content
	} else {
		rows, err := db.Query(`SELECT json FROM lines WHERE file_id = ? ORDER BY n`, id)
		if err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		buffer.Grow(int(size))
		first := true
		for rows.Next() {
			var line string
			if err := rows.Scan(&line); err != nil {
				rows.Close()
				return nil, err
			}
			if !first {
				buffer.WriteByte('\n')
			}
			buffer.WriteString(line)
			first = false
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if !first && !tailOpen {
			buffer.WriteByte('\n')
		}
		file.Data = buffer.Bytes()
	}

	if int64(len(file.Data)) != size || hashOf(file.Data) != hash {
		return nil, errors.New("stored copy is damaged")
	}
	return &file, nil
}

func idsFor(path string) (project string, session string) {
	parts := strings.Split(path, "/")
	switch parts[0] {
	case "projects":
		if len(parts) >= 2 {
			project = parts[1]
		}
		if len(parts) >= 3 {
			candidate := parts[2]
			if len(parts) == 3 {
				candidate, _, _ = strings.Cut(candidate, ".")
			}
			if uuidPattern.MatchString(candidate) {
				session = candidate
			}
		}
	case "file-history", "tasks":
		if len(parts) >= 2 && uuidPattern.MatchString(parts[1]) {
			session = parts[1]
		}
	}
	return project, session
}

func validPath(path string) bool {
	return path != "" && !strings.Contains(path, `\`) && !strings.HasPrefix(path, "/") && filepath.IsLocal(path)
}

func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func hashOf(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
