package tcphandle

import (
	"fmt"
	"net"
	"time"
)

type Service struct {
	slaves map[string]*Slave
}

var (
	writeTimeout, readTimeout = 5 * time.Second, 5 * time.Second
	service                   *Service
)

func init() {
	service = &Service{}
	service.slaves = make(map[string]*Slave)
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
		fmt.Println("accept conn:", string(versionbytes))
		slave := &Slave{service: s}
		handler, err := newVersionHandler(versionbytes, slave)
		if err != nil {
			fmt.Println("unsupport protocol version", string(versionbytes))
			continue
		}
		key := conn.(*net.TCPConn).RemoteAddr().String()
		if s.slaves == nil {
			s.slaves = make(map[string]*Slave)
		}
		s.slaves[key] = slave
		slave.handle = handler
		go slave.handle.MainLoop(conn)
	}
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

type ProtocolHandler interface {
	MainLoop(net.Conn)
}
