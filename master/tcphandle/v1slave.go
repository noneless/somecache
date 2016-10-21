package tcphandle

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type V1Slave struct {
	conn   net.Conn
	reader bufio.Reader
	lock   sync.Mutex
	ctx    context
}

func (v1s *V1Slave) CommandLoop() {
	pingticker := time.NewTicker(time.Second)
	for {
		select {
		case <-pingticker.C:
			v1s.lock.Lock()
			v1s.lock.Unlock()

		}
	}
}

func (v1s *V1Slave) Ping() error {
	return nil
}
