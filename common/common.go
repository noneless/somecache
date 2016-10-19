package common

import (
	"net"
)

type Getter interface {
	Get(string) []byte
}

type TcpServer interface {
	TcpServer(net.Listener)
}

var (
	MagicV1     = []byte("  V1")
	COMMAND_PUT = []byte("PUT")
	COMMAND_GET = []byte("GET")
	E_ERROR     = []byte("E_ERROR")
	WhiteSpace  = []byte(' ')
)

type Command struct {
}
