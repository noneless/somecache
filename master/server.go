package master

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/consistenthash"
)

type Service struct {
	hosts map[string]*Host
	lock  sync.Mutex
	hash  *consistenthash.Map
}

var (
	writeTimeout, readTimeout = 5 * time.Second, 5 * time.Second
	service                   *Service
	defaultReplicas           = 50
)

func init() {
	service = &Service{}
	service.hosts = make(map[string]*Host)
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
		index := strings.LastIndex(key, ":")
		if s.hosts == nil {
			s.hosts = make(map[string]*Host)
		}
		slave.handle = handler
		c := make(chan struct{})
		go func() {
			slave.handle.MainLoop(conn, c)
			s.delSlave(key[0:index], key[index+1:])
		}()
		select {
		case <-c: // slave is set up ok,ready to service
		}
		s.addSlave(key[0:index], key[index+1:], slave)
	}
}

func (s *Service) addSlave(host string, port string, slave *Slave) {
	s.lock.Lock()
	defer s.lock.Unlock()
	h, ok := s.hosts[host]
	if !ok {
		h = &Host{
			workers: make(map[string]*Slave),
		}
		s.hosts[host] = h
	}
	h.addSlave(port, slave)
	s.reBuildHash()
}

func (s *Service) reBuildHash() {
	s.hash.Empty()
	keys := make([]string, 0)
	for k, _ := range s.hosts {
		keys = append(keys, k)
	}
	s.hash.Add(keys...)
}

func (s *Service) delSlave(host, port string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	h, ok := s.hosts[host]
	if ok {
		h.delSlave(port)
	}
	s.reBuildHash()
}

//get slave is going to get a download worker
func (s *Service) getSlave(key string) *Host {
	s.lock.Lock()
	defer s.lock.Unlock()
	key = s.hash.Get(key)
	return s.hosts[key]
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

type ProtocolHandler interface {
	MainLoop(net.Conn, chan struct{})
	Close()
	Exec(c *common.Command, fn func(io.Reader, int) error)
	IfBusy() int64
}
