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
	"bufio"
	"flag"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/756445638/somecache/common"
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
	go func() {
		err = master.Server(ln)
		if err != nil {
			fmt.Printf("master.Server server failed,err[%v]\n", err)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	var s string
	var err error
	for err == nil {

	}

	wg.Add(1)

}
