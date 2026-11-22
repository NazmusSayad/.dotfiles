package agentstore

import (
	"bufio"
	"encoding/gob"
	"io"
)

const ProtocolVersion = 1

const ReadyPrefix = "AGENTS-SERVER READY"

const (
	KindLines = "lines"
	KindBlob  = "blob"
)

const (
	ModeFull   = "full"
	ModeAppend = "append"
	ModeTouch  = "touch"
)

type Hello struct {
	Version int
	Token   string
	Op      string
	Device  string
}

type FileState struct {
	Path  string
	Kind  string
	Size  int64
	Hash  string
	Mtime int64
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
	Base  int64
	Size  int64
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

type StatsGroup struct {
	Device     string
	Model      string
	Month      string
	Requests   int64
	Input      int64
	Output     int64
	CacheWrite int64
	CacheRead  int64
}

type Stats struct {
	Groups          []StatsGroup
	Sessions        int
	Storage         string
	UnreadableLines int
	Err             string
}

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
