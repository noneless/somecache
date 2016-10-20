package main

import (
	"flag"
	"net"
	"sync"

	//	"github.com/756445638/somecache/master/grpc"
)

var (
	wg           = sync.WaitGroup{}
	tcp_address  = flag.Int("tcp-address", 4000, "tcp address")
	http_address = flag.Int("http-address", 4001, "http server address")
)

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", *tcp_address)
}
