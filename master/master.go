package main

import (
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/756445638/somecache/master/tcphandle"
)

var (
	wg          = sync.WaitGroup{}
	tcp_address = flag.Int("tcp-address", 4000, "tcp address")
)

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *tcp_address))
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go func() {
		err := tcphandle.Server(ln)
		wg.Done()
		panic(err)
	}()
	wg.Wait()

}
