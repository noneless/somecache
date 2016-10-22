package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/756445638/somecache/master"
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
		defer wg.Done()
		err := master.Server(ln)
		if err != nil {
			fmt.Printf("master.Server server failed,err[%v]\n", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fileserver()
	}()
	wg.Wait()
}

func fileserver() {

}
