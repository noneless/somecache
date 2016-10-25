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
