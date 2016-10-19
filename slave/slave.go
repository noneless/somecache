package main

import (
	"flag"
	"net"
	"sync"
)

var (
	master = flag.String("master", "", "master addr")
	tcp    = flag.String("tcp-address", "", "tcp address")
	wg     = sync.WaitGroup{}
)

func main() {
	flag.Parse()
	if *master == "" {
		panic("master is nil string")
	}
	if *tcp == "" {
		panic("tcp is nil string")
	}
	ln, err := net.Listen("tcp", *tcp)
	if err != nil {
		panic(err)
	}
	go HandleTcp(ln)
	go handlemaster()
	wg.Wait()
}
