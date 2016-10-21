package tcphandle

import (
	"bytes"
	"fmt"
	"net"

	"github.com/756445638/somecache/common"
)

type Service struct {
	slaves map[string]SlaveHandler
}

var service *Service

func init() {
	service = &Service{}
	service.slaves = make(map[string]SlaveHandler)
}

func (s *Service) Server(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		versionbytes := make([]byte, 4)
		n, err := conn.Read(versionbytes)
		if err != nil || n != 4 {
			fmt.Println("accept error:", err)
			continue
		}
		if bytes.Equal(versionbytes, common.MagicV1) {
			v1s := &V1Slave{conn: conn}
			v1s.ctx.service = s
			key := conn.(*net.TCPConn).RemoteAddr().String()
			service.slaves[key] = v1s
			go func() {
				service.slaves[key].CommandLoop()
			}()
		} else {
			fmt.Println("unsupport protocol version")
		}
	}
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

type SlaveHandler interface {
	CommandLoop()
}