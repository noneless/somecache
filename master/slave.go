package master

import (
	"bytes"
	"errors"
	"io"
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
			jobschan:  make(chan *job),
		}
		return v1, nil
	} else {
		return nil, errors.New("unkown version")
	}
}

type ProtocolHandler interface {
	MainLoop(net.Conn, chan bool) //chan bool means if this woker is setup ok
	Close()
	Get2Stream(key string, w io.Writer) error // stream way to get cache
	Get(key string) ([]byte, error)           // read it to memory
	Put(key string, data []byte) error
	PutFromReader(key string, reader io.Reader) error
	IfBusy() int64
}
