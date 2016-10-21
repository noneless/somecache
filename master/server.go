package master

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/756445638/somecache/consistenthash"
)

type Service struct {
	slaves map[string]*Slave
	lock   sync.Mutex
	hash   *consistenthash.Map
}

var (
	writeTimeout, readTimeout = 5 * time.Second, 5 * time.Second
	service                   *Service
	defaultReplicas           = 50
)

func init() {
	service = &Service{}
	service.slaves = make(map[string]*Slave)
	service.hash = consistenthash.New(defaultReplicas, nil)
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
			slave.handle.MainLoop(conn, c)
			s.delSlave(key)
		}()
		select {
		case <-c: // slave is set up ok,ready to service
		}
		s.addSlave(key, slave)
	}
}

func (s *Service) addSlave(key string, slave *Slave) {
	s.lock.Lock()
	defer s.lock.Unlock()
	e, ok := s.slaves[key]
	if ok {
		e.handle.Close()
	}
	s.slaves[key] = slave
	s.reBuildHash()
}

func (s *Service) reBuildHash() {
	s.hash.Empty()
	keys := make([]string, 0)
	for k, _ := range s.slaves {
		keys = append(keys, k)
	}
	s.hash.Add(keys...)
}

func (s *Service) delSlave(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.slaves, key)
	s.reBuildHash()
}

func (s *Service) getSlave(key string) *Slave {
	s.lock.Lock()
	defer s.lock.Unlock()
	key = s.hash.Get(key)
	return s.slaves[key]
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

type ProtocolHandler interface {
	MainLoop(net.Conn, chan struct{})
	Close()
}
