package master

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
			fmt.Println("read version bytes error:", err)
			continue
		}
		fmt.Println("accept conn:", string(versionbytes), time.Now())
		slave := &Slave{service: s}
		handler, err := newVersionHandler(versionbytes, slave)
		if err != nil {
			fmt.Println("unsupport protocol version", string(versionbytes))
			continue
		}
		slave.addr = conn.(*net.TCPConn).RemoteAddr()
		key := slave.addr.String()
		if s.slaves == nil {
			s.slaves = make(map[string]*Slave)
		}
		slave.handle = handler
		c := make(chan struct{})
		go func() {
			go slave.handle.MainLoop(conn, c)
			select {
			case <-c: // slave is set up ok,ready to service
			}
			s.AddSlave(key, slave)
		}()
	}
}

func (s *Service) AddSlave(key string, slave *Slave) {
	s.slaves[key] = slave
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

type ProtocolHandler interface {
	MainLoop(net.Conn, chan struct{})
	Close()
}
