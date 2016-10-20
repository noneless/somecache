package common

import (
	"net"
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

func NewCommand(command []byte, paras [][]byte, content []byte) *Command {
	return &Command{command, paras, content}
}
