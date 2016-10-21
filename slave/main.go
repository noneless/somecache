package main

import (
	"flag"
	"net"
	"sync"
)

var (
	master = flag.String("master-tcp-address", "127.0.0.1:4000", "master addr")
	tcp    = flag.String("tcp-address", ":4001", "tcp address")
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
	wg.Add(2)
	go HandleTcp(ln)
	go Connection2Master(*master)
	wg.Wait()
}
