package common

import (
	"bytes"
	"io"
	"net"
	//	"github.com/756445638/somecache/lru"
)

type Getter interface {
	Get(string) ([]byte, error)
}

type TcpServer interface {
	TcpServer(net.Listener)
}

var (
	MagicV1     = []byte("  V1")
	COMMAND_PUT = []byte("PUT")
	COMMAND_GET = []byte("GET")
	E_ERROR     = []byte("E_ERROR")
	E_NOT_FOUND = []byte("NOT_FOUND")
	ENDL        = []byte("\n")
	OK          = []byte("OK")
	WhiteSpace  = []byte(" ")
)

type Command struct {
	Command   []byte
	Parameter [][]byte
	Content   []byte
}

func (c *Command) Write(w io.Writer) (int, error) {
	total := int(0)
	n, err := w.Write(c.Command)
	total += n
	if err != nil {
		return total, err
	}
	n, err = w.Write(bytes.Join(c.Parameter, WhiteSpace))
	total += n
	if err != nil {
		return total, err
	}
	n, err = w.Write(ENDL)
	total += n
	if err != nil {
		return total, err
	}

	if c.Content != nil && len(c.Content) > 0 {
		length := len(c.Content)
		for length > 0 {
			n, err = w.Write(c.Content)
			total += n
			if err != nil {
				return total, err
			}
			length -= n
			c.Content = c.Content[n:]
		}
	}
	return total, nil
}

//not thread safe
func NewCommand(command []byte, paras [][]byte, content []byte) *Command {
	return &Command{command, paras, content}
}

type BytesDate []byte

func (bd BytesDate) Measure() uint64 {
	return uint64(len(bd))
}
