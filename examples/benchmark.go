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
package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/756445638/somecache/master"
)

var (
	wg          = sync.WaitGroup{}
	tcp_address = flag.Int("tcp-address", 4000, "tcp address")
	dir         = flag.String("dir", "", "directory")
)

func signalHandle() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)
	x := <-c
	panic(x.String())
}

func main() {
	go signalHandle()
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *tcp_address))
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = master.Server(ln)
		if err != nil {
			fmt.Printf("master.Server server failed,err[%v]\n", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 5)

		for i := 0; i < 5; i++ {
			runBenchMark()
		}
	}()
	wg.Wait()
}

func runBenchMark() {
	s := 1024 * 50
	var total time.Duration
	for i := 0; i < s; i++ {
		now := time.Now()
		k := fmt.Sprintf("k%d", i)
		_, err := master.Get(k)
		if err != nil {
			fmt.Println("get error:", err)
		}
		total += time.Now().Sub(now)
	}
	fmt.Printf("get %d objects,takes: \n", s, total.Seconds())

}

func init() {
	master.RegisterGetter(&MemoryGetter{})
}

type MemoryGetter struct {
}

func (MemoryGetter) Get(k string) ([]byte, error) {
	length := crc32.ChecksumIEEE([]byte(k))
	length = length % (4 * 1024)
	s := time.Now().String() + k
	for uint32(len(s)) < length { // less than 4 * 1024
		s += s
	}
	return []byte(s), nil
}
