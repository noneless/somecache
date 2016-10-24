package master

import (
	"bytes"
	"errors"
	"net"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/message"
)

type Slave struct {
	addr         net.Addr
	service      *Service
	handle       ProtocolHandler
	loginmessage *message.Login
}

func newVersionHandler(v []byte, slave *Slave) (ProtocolHandler, error) {
	if bytes.Equal(v, common.MagicV1) {
		v1 := &V1Slave{
			slave:     slave,
			closechan: make(chan struct{}),
			notify:    make(map[uint64]*job),
			jobschan:  make(chan *job, 1024),
		}
		return v1, nil
	} else {
		return nil, errors.New("unkown version")
	}
}

type ProtocolHandler interface {
	MainLoop(net.Conn, chan bool) //chan bool means if this woker is setup ok
	Close()
	Get(key string) ([]byte, error) // read it to memory
	Put(key string, data []byte) error
}
