/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package master

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/consistenthash"
	"github.com/756445638/somecache/lru"
)

type Service struct {
	hosts      map[string]*Host
	lock       sync.Mutex
	hash       *consistenthash.Map
	localcache lru.Lru
	jobchan    chan *job
}

var (
	writeTimeout, readTimeout = 30 * time.Second, 30 * time.Second
	service                   *Service
	defaultReplicas                 = 50
	defaultCacheSize          int64 = 1 << 30
)

func init() {
	service = &Service{}
	service.hosts = make(map[string]*Host)
	service.hash = consistenthash.New(defaultReplicas, nil)
}

func SetUpCacheSize(size int64) {
	defaultCacheSize = size
}

func (s *Service) setDefaultParameter() {
	s.localcache.SetMaxCacheSize(defaultCacheSize)
	s.jobchan = make(chan *job, 1024)
}

func (s *Service) Server(ln net.Listener) error {
	s.setDefaultParameter()
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

func (s *Service) getLocalCache(key string) []byte {
	v := s.localcache.Get(key)
	if v == nil {
		return nil
	}
	data := v.(*common.BytesData)
	return data.Data
}
func (s *Service) putLocalCache(key string, d []byte) {
	s.localcache.Put(key, &common.BytesData{K: key, Data: d})

}

func (s *Service) getRemoteCache(key string) ([]byte, error) {
	worker := s.getSlave(key)
	if worker == nil {
		return nil, fmt.Errorf("no worker available")
	}
	return worker.handle.Get(key)
}

func (s *Service) putRemoteCache(key string, d []byte) error {
	worker := s.getSlave(key)
	if worker == nil {
		return fmt.Errorf("no worker available")
	}
	return worker.handle.Put(key, d)
}

// main entrance
func (s *Service) Get(key string) ([]byte, error) {
	data := s.getLocalCache(key)
	if data != nil { // find cache in localcache,nothing to do,just return
		return data, nil
	}
	var err error
	data, err = s.getRemoteCache(key)
	if err == nil { // find cache in remote cache,store in local
		s.putLocalCache(key, data)
		return data, nil
	}
	// not in localcache and not in remotecache

	if !strings.Contains(err.Error(), "NOT_FOUND") && !strings.Contains(err.Error(), "timeout") {
		//some error but no "NOT_FOUND",it is a not very serious error
		return nil, err
	}
	fmt.Println("remote server error:", err)
	if getter == nil { //getter is not registed
		return nil, fmt.Errorf("not found  and getter is not registered")
	}
	data, err = getter.Get(key) // call getter returned error
	if err != nil {
		return nil, err
	}
	s.putLocalCache(key, data)
	go func() {
		err := s.putRemoteCache(key, data)
		fmt.Println("pus in remote server error:", err)
	}()
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
