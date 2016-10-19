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

const (
	MagicV1 = []byte("  V1")
)
