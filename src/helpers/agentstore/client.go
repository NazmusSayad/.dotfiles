package agentstore

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

var claudeItems = []string{
	"projects",
	"history.jsonl",
	"file-history",
	"plans",
	"tasks",
	"paste-cache",
}

var pullSkip = map[string]bool{
	"history.jsonl": true,
}

const tempPrefix = ".agents-pull-"

type localFile struct {
	rel   string
	abs   string
	size  int64
	mtime int64
}

type link struct {
	conn  *Conn
	token string
	stdin io.WriteCloser
	cmds  []*exec.Cmd
}

func fatal(message string) {
	fmt.Println(aurora.Red("Error: " + message))
	os.Exit(1)
}

func claudeDir() string {
	if dir := os.Getenv("AGENTS_CLAUDE_DIR"); dir != "" {
		return helpers.ResolvePath(dir)
	}
	return helpers.ResolvePath("~/.claude")
}

func serverCommand(device string) []string {
	if command := os.Getenv("AGENTS_SERVER"); command != "" {
		return strings.Fields(command)
	}
	if device == "wsl" {
		return []string{"agents-server.exe"}
	}
	return []string{"agents-server"}
}

func connect(device string) *link {
	command := serverCommand(device)
	if device == helpers.CLAUDE_HOST_DEVICE || device == "wsl" {
		return startLocal(command)
	}
	return startRemote(command)
}

func startLocal(command []string) *link {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		fatal("Failed to start " + command[0] + ": " + err.Error())
	}
	return &link{conn: NewConn(stdout, stdin), stdin: stdin, cmds: []*exec.Cmd{cmd}}
}

func startRemote(command []string) *link {
	sshEnv := helpers.CLAUDE_DEVICES[helpers.CLAUDE_HOST_DEVICE].SSHEnv
	address := os.Getenv(sshEnv)
	if address == "" {
		fatal(sshEnv + " must be set")
	}

	server := exec.Command("ssh", "-n", "-o", "BatchMode=yes", address, strings.Join(command, " ")+" --listen")
	server.Stderr = os.Stderr
	serverOut, err := server.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := server.Start(); err != nil {
		fatal("Failed to start ssh: " + err.Error())
	}

	ready := make(chan []string, 1)
	go func() {
		scanner := bufio.NewScanner(serverOut)
		for scanner.Scan() {
			line := scanner.Text()
			if index := strings.Index(line, ReadyPrefix); index >= 0 {
				ready <- strings.Fields(line[index+len(ReadyPrefix):])
				break
			}
		}
		close(ready)
		_, _ = io.Copy(io.Discard, serverOut)
	}()

	var fields []string
	select {
	case fields = <-ready:
	case <-time.After(3 * time.Minute):
		_ = server.Process.Kill()
		fatal("agents-server on " + address + " did not start in time")
	}
	if len(fields) != 2 {
		_ = server.Wait()
		fatal("agents-server on " + address + " did not start")
	}

	tunnel := exec.Command("ssh", "-o", "BatchMode=yes", "-W", "127.0.0.1:"+fields[0], address)
	tunnel.Stderr = os.Stderr
	stdin, err := tunnel.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := tunnel.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := tunnel.Start(); err != nil {
		_ = server.Process.Kill()
		fatal("Failed to start ssh tunnel: " + err.Error())
	}

	return &link{conn: NewConn(stdout, stdin), token: fields[1], stdin: stdin, cmds: []*exec.Cmd{tunnel, server}}
}

func (l *link) close() {
	_ = l.stdin.Close()
	for _, cmd := range l.cmds {
		_ = cmd.Wait()
	}
}

func (l *link) sendHello(op string, device string) {
	if err := l.conn.Send(Hello{Version: ProtocolVersion, Token: l.token, Op: op, Device: device}); err != nil {
		fatal("Failed to reach agents-server: " + err.Error())
	}
}

func (l *link) hello(op string, device string) State {
	l.sendHello(op, device)
	var state State
	if err := l.conn.Receive(&state); err != nil {
		fatal("Failed to read from agents-server: " + err.Error())
	}
	if state.Err != "" {
		fatal(state.Err)
	}
	return state
}

func kindOf(rel string) string {
	if strings.HasSuffix(rel, ".jsonl") {
		return KindLines
	}
	return KindBlob
}

func hashOf(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func listLocal(root string) []localFile {
	var files []localFile
	for _, item := range claudeItems {
		err := filepath.WalkDir(filepath.Join(root, item), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			}
			if !entry.Type().IsRegular() {
				return nil
			}
			name := entry.Name()
			if strings.HasSuffix(name, ".wakatime") || strings.HasPrefix(name, tempPrefix) {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, localFile{rel: filepath.ToSlash(rel), abs: path, size: info.Size(), mtime: info.ModTime().UnixMilli()})
			return nil
		})
		if err != nil {
			fatal("Failed to read " + filepath.Join(root, item) + ": " + err.Error())
		}
	}
	return files
}

func planUpload(file localFile, stored *FileState) (*Upload, error) {
	if stored != nil && stored.Size == file.size && stored.Mtime == file.mtime {
		return nil, nil
	}

	data, err := os.ReadFile(file.abs)
	if err != nil {
		return nil, err
	}
	upload := &Upload{Path: file.rel, Kind: kindOf(file.rel), Size: int64(len(data)), Hash: hashOf(data), Mtime: file.mtime}

	switch {
	case stored == nil || stored.Kind != upload.Kind:
		upload.Mode = ModeFull
		upload.Data = data
	case stored.Hash == upload.Hash:
		upload.Mode = ModeTouch
	case upload.Kind == KindLines && upload.Size > stored.Size && hashOf(data[:stored.Size]) == stored.Hash:
		upload.Mode = ModeAppend
		upload.Base = stored.Size
		upload.Data = data[stored.Size:]
	default:
		upload.Mode = ModeFull
		upload.Data = data
	}
	return upload, nil
}

func Push() {
	started := time.Now()
	device := helpers.DetectClaudeDevice()
	root := claudeDir()
	if _, err := os.Stat(root); err != nil {
		fatal("No Claude data found in " + root)
	}

	fmt.Println(aurora.Faint("Reading " + root))
	local := listLocal(root)
	if len(local) == 0 {
		fatal("No Claude data found in " + root)
	}
	var localBytes int64
	for _, file := range local {
		localBytes += file.size
	}
	fmt.Printf("Found %d files (%s)\n", len(local), formatBytes(localBytes))

	fmt.Println(aurora.Faint("Connecting to agents-server"))
	l := connect(device)
	defer l.close()
	state := l.hello("push", device)

	stored := make(map[string]*FileState, len(state.Files))
	for i := range state.Files {
		stored[state.Files[i].Path] = &state.Files[i]
	}

	var changed []localFile
	var changedBytes int64
	for _, file := range local {
		if known := stored[file.rel]; known == nil || known.Size != file.size || known.Mtime != file.mtime {
			changed = append(changed, file)
			changedBytes += file.size
		}
	}
	fmt.Printf("Stored for %s: %d files. Changed since then: %d files (%s)\n", device, len(state.Files), len(changed), formatBytes(changedBytes))

	var readErrors []string
	var sentBytes int64
	width := len(strconv.Itoa(len(changed)))
	for i, file := range changed {
		upload, err := planUpload(file, stored[file.rel])
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				readErrors = append(readErrors, file.rel+": "+err.Error())
			}
			continue
		}
		if upload == nil {
			continue
		}
		if err := l.conn.Send(upload); err != nil {
			fatal("Failed to send " + file.rel + ": " + err.Error())
		}
		sentBytes += int64(len(upload.Data))

		var label aurora.Value
		switch {
		case upload.Mode == ModeTouch:
			label = aurora.Faint("time    ")
		case upload.Mode == ModeAppend:
			label = aurora.Blue("grown   ")
		case stored[file.rel] == nil:
			label = aurora.Green("new     ")
		default:
			label = aurora.Yellow("replaced")
		}
		fmt.Printf("%s %s %s %s\n",
			aurora.Faint(fmt.Sprintf("[%*d/%d %s]", width, i+1, len(changed), formatBytes(sentBytes))),
			label, file.rel, aurora.Faint(formatBytes(int64(len(upload.Data)))))
	}

	if err := l.conn.Send(Upload{Done: true}); err != nil {
		fatal("Failed to finish push: " + err.Error())
	}
	if len(changed) > 0 {
		fmt.Println(aurora.Faint("Waiting for agents-server to save"))
	}
	var result PushResult
	if err := l.conn.Receive(&result); err != nil {
		fatal("Failed to read push result: " + err.Error())
	}
	if result.Err != "" {
		fatal(result.Err)
	}

	for _, message := range append(readErrors, result.Errors...) {
		fmt.Println(aurora.Red("Skipped " + message))
	}
	fmt.Println(aurora.Green(fmt.Sprintf(
		"Pushed %s: %d new, %d grown, %d replaced, %d with only a new time (%s in %s)",
		device, result.Added, result.Appended, result.Replaced, result.Touched, formatBytes(result.Bytes), formatDuration(time.Since(started)),
	)))
}

func Pull() {
	started := time.Now()
	device := helpers.DetectClaudeDevice()
	root := claudeDir()

	fmt.Println(aurora.Faint("Connecting to agents-server"))
	l := connect(device)
	defer l.close()
	state := l.hello("pull", device)

	var wants []string
	var wantBytes int64
	for _, stored := range state.Files {
		if pullSkip[stored.Path] {
			continue
		}
		rel := filepath.FromSlash(stored.Path)
		if !filepath.IsLocal(rel) {
			continue
		}
		info, err := os.Lstat(filepath.Join(root, rel))
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			continue
		case stored.Kind == KindLines && info.Mode().IsRegular() && info.Size() < stored.Size:
		default:
			continue
		}
		wants = append(wants, stored.Path)
		wantBytes += stored.Size
	}
	fmt.Printf("Stored for %s: %d files. Missing or shorter in %s: %d files (%s)\n", device, len(state.Files), root, len(wants), formatBytes(wantBytes))

	if err := l.conn.Send(Want{Paths: wants}); err != nil {
		fatal("Failed to send pull request: " + err.Error())
	}

	var restored, extended, differ, failed int
	var receivedBytes int64
	width := len(strconv.Itoa(len(wants)))
	for count := 1; ; count++ {
		var file File
		if err := l.conn.Receive(&file); err != nil {
			fatal("Failed to read from agents-server: " + err.Error())
		}
		if file.Err != "" {
			fatal(file.Err)
		}
		if file.Done {
			break
		}
		receivedBytes += int64(len(file.Data))

		outcome, err := applyFile(root, file)
		var label aurora.Value
		switch {
		case err != nil:
			failed++
			label = aurora.Red("failed  ")
		case outcome == "restored":
			restored++
			label = aurora.Green("restored")
		case outcome == "extended":
			extended++
			label = aurora.Blue("extended")
		case outcome == "differ":
			differ++
			label = aurora.Yellow("kept    ")
		default:
			label = aurora.Faint("current ")
		}
		line := fmt.Sprintf("%s %s %s %s",
			aurora.Faint(fmt.Sprintf("[%*d/%d %s]", width, count, len(wants), formatBytes(receivedBytes))),
			label, file.Path, aurora.Faint(formatBytes(int64(len(file.Data)))))
		switch {
		case err != nil:
			line += " " + aurora.Red(err.Error()).String()
		case outcome == "differ":
			line += " " + aurora.Yellow("local copy differs, left as is").String()
		}
		fmt.Println(line)
	}

	summary := fmt.Sprintf("Pulled %s into %s: %d restored, %d extended", device, root, restored, extended)
	if differ > 0 {
		summary += fmt.Sprintf(", %d kept because the local copy differs", differ)
	}
	if failed > 0 {
		summary += fmt.Sprintf(", %d failed", failed)
	}
	summary += " (" + formatDuration(time.Since(started)) + ")"
	if failed > 0 {
		fmt.Println(aurora.Red(summary))
	} else {
		fmt.Println(aurora.Green(summary))
	}
}

func applyFile(root string, file File) (string, error) {
	rel := filepath.FromSlash(file.Path)
	if !filepath.IsLocal(rel) {
		return "", fmt.Errorf("unsafe path")
	}
	abs := filepath.Join(root, rel)
	mtime := time.UnixMilli(file.Mtime)

	local, err := os.ReadFile(abs)
	if errors.Is(err, fs.ErrNotExist) {
		return "restored", writeNew(abs, file.Data, mtime)
	}
	if err != nil {
		return "", err
	}

	if len(local) >= len(file.Data) {
		return "", nil
	}
	if file.Kind != KindLines || !bytes.HasPrefix(file.Data, local) {
		return "differ", nil
	}

	handle, err := os.OpenFile(abs, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return "", err
	}
	info, err := handle.Stat()
	if err != nil {
		handle.Close()
		return "", err
	}
	if info.Size() != int64(len(local)) {
		handle.Close()
		return "", fmt.Errorf("changed while pulling, left as is")
	}
	if _, err := handle.Write(file.Data[len(local):]); err != nil {
		handle.Close()
		return "", err
	}
	if err := handle.Close(); err != nil {
		return "", err
	}
	return "extended", os.Chtimes(abs, mtime, mtime)
}

func writeNew(abs string, data []byte, mtime time.Time) error {
	dir := filepath.Dir(abs)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(dir, tempPrefix+"*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Chtimes(tempPath, mtime, mtime); err != nil {
		return err
	}
	if _, err := os.Lstat(abs); err == nil {
		return fmt.Errorf("appeared while pulling, left as is")
	}
	return os.Rename(tempPath, abs)
}

func formatBytes(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(100 * time.Millisecond).String()
	}
	return d.Round(time.Second).String()
}

func FetchStats() Stats {
	device := helpers.DetectClaudeDevice()
	l := connect(device)
	defer l.close()
	l.sendHello("stats", device)

	var stats Stats
	if err := l.conn.Receive(&stats); err != nil {
		fatal("Failed to read from agents-server: " + err.Error())
	}
	if stats.Err != "" {
		fatal(stats.Err)
	}
	return stats
}
