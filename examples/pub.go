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
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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
	for err == nil {
		s, err = reader.ReadString('\n')
		s = strings.TrimRight(s, "\n")
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "get ") {
			s = strings.TrimLeft(s, "get ")
			s = strings.TrimLeft(s, " ")
			get(s)
		} else if strings.HasPrefix(s, "put ") {
			s = strings.TrimLeft(s, "put ")
			s = strings.TrimLeft(s, " ")
			put(s)
		} else {
			fmt.Println("unkown command")
		}
	}
	wg.Add(1)
}

func get(s string) {
	remote := false
	if strings.HasPrefix(s, "remote ") {
		s = strings.TrimLeft(s, "remote ")
		remote = true
	}
	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")
	var data []byte
	var err error
	if remote { //get remote just for test ,in porduction alway look localcache first
		data, err = master.GetFromRemoteServer(s)
	} else {
		data, err = master.Get(s)
	}
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(string(data))
}

func put(s string) {
	remote := false
	if strings.HasPrefix(s, "remote ") {
		s = strings.TrimLeft(s, "remote ")
		remote = true
	}
	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")
	t := strings.Split(s, " ")
	if len(t) < 2 {
		fmt.Println("paramter error")
		return
	}

	key := t[0]
	data := []byte(t[1])
	master.Put(key, data, remote)

}
