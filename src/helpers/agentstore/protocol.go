// Package agentstore moves Claude data between a device and the agents
// database on the Windows host. The database itself is only opened by
// agents-server; devices talk to it through this protocol.
package agentstore

import (
	"bufio"
	"encoding/gob"
	"io"
)

const ProtocolVersion = 1

// ReadyPrefix is printed by agents-server in listen mode, followed by the
// port and a one-time token.
const ReadyPrefix = "AGENTS-SERVER READY"

const (
	KindLines = "lines" // .jsonl files, stored one JSON line per row
	KindBlob  = "blob"  // everything else, stored as raw bytes
)

const (
	ModeFull   = "full"   // store the whole file, replacing what is there
	ModeAppend = "append" // add bytes to the end of the stored file
	ModeTouch  = "touch"  // content is unchanged, only update the modified time
)

type Hello struct {
	Version int
	Token   string
	Op      string // "push" or "pull"
	Device  string
}

type FileState struct {
	Path  string
	Kind  string
	Size  int64
	Hash  string
	Mtime int64 // unix milliseconds
}

type State struct {
	Files []FileState
	Err   string
}

type Upload struct {
	Done  bool
	Path  string
	Kind  string
	Mode  string
	Base  int64 // stored size the append starts from
	Size  int64 // full size after this upload
	Hash  string
	Mtime int64
	Data  []byte
}

type PushResult struct {
	Added    int
	Appended int
	Replaced int
	Touched  int
	Bytes    int64
	Errors   []string
	Err      string
}

type Want struct {
	Paths []string
}

type File struct {
	Done  bool
	Path  string
	Kind  string
	Mtime int64
	Data  []byte
	Err   string
}

// Conn sends and receives protocol messages over any byte stream.
type Conn struct {
	w   *bufio.Writer
	enc *gob.Encoder
	dec *gob.Decoder
}

func NewConn(r io.Reader, w io.Writer) *Conn {
	bw := bufio.NewWriterSize(w, 1<<20)
	return &Conn{w: bw, enc: gob.NewEncoder(bw), dec: gob.NewDecoder(bufio.NewReaderSize(r, 1<<20))}
}

func (c *Conn) Send(v any) error {
	if err := c.enc.Encode(v); err != nil {
		return err
	}
	return c.w.Flush()
}

func (c *Conn) Receive(v any) error {
	return c.dec.Decode(v)
}
