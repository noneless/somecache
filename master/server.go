package master

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/756445638/somecache/consistenthash"
	"github.com/756445638/somecache/lru"
)

type Service struct {
	hosts      map[string]*Host
	lock       sync.Mutex
	hash       *consistenthash.Map
	localcache lru.Lru
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
		c := make(chan bool)
		setupok := false
		go func(sok *bool) {
			slave.handle.MainLoop(conn, c)
			if *sok {
				s.delSlave(key[0:index], key[index+1:])
			}
		}(&setupok)
		select {
		case setupok = <-c: // slave is set up ok,ready to server
		}
		if setupok {
			s.addSlave(key[0:index], key[index+1:], slave)
		}
	}
}

func (s *Service) addSlave(host string, port string, slave *Slave) {
	s.lock.Lock()
	defer s.lock.Unlock()
	h, ok := s.hosts[host]
	rebuild := false
	if !ok {
		h = &Host{
			workers: make(map[string]*Slave),
		}
		s.hosts[host] = h
		rebuild = true
	}
	h.addSlave(port, slave)
	if rebuild {
		s.reBuildHash()
	}
}

func (s *Service) Get(key string) ([]byte, error) {
	worker := s.getSlave(key)
	if worker == nil {
		return nil, fmt.Errorf("no worker available")
	}
	b, err := worker.handle.Get(key)
	if err == nil {
		return b, nil
	}
	if !strings.Contains(err.Error(), "NOT_FOUND") {
		return nil, err
	}
	if getter == nil {
		return nil, fmt.Errorf("not found  and getter is not registered")
	}
	data, err := getter.Get(key)
	if err != nil {
		return nil, err
	}
	return data, nil
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
	rebuild := false
	if ok {
		h.delSlave(port)
		if len(h.workers) == 0 { // no more worker
			rebuild = true
			delete(s.hosts, host)
		}
	}
	if rebuild {
		s.reBuildHash()
	}
}

//get slave is going to get a download worker
func (s *Service) getSlave(key string) *Slave {
	s.lock.Lock()
	defer s.lock.Unlock()
	key = s.hash.Get(key)
	h, ok := s.hosts[key]
	if !ok {
		return nil
	}
	return h.getWorker() // get a download worker
}

func Server(ln net.Listener) error {
	return service.Server(ln)
}

func validKey(key string) bool {
	return !strings.Contains(key, "\n")
}

var InValidKeyError = errors.New("invalid name error")
