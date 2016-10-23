package master

import (
	"net"

	"github.com/756445638/somecache/message"
)

type Slave struct {
	addr         net.Addr
	service      *Service
	handle       ProtocolHandler
	loginmessage *message.Login
}
