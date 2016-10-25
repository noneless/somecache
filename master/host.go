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
	"sync"
	"time"
)

type Host struct {
	lock    sync.Mutex
	host    string
	workers map[string]*Slave
}

func (h *Host) addSlave(port string, s *Slave) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.workers[port] = s
}

func (h *Host) delSlave(port string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.workers, port)
}

func (h *Host) getWorker() *Slave {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.workers) == 0 {
		return nil
	}
	index := time.Now().Nanosecond() % len(h.workers)
	i := 0
	for _, v := range h.workers {
		if i == index {
			return v
		}
		i++
	}
	return nil
}
