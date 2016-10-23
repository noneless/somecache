package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"net"
	"sync"
	"time"

	"github.com/756445638/somecache/master"
)

var (
	wg          = sync.WaitGroup{}
	tcp_address = flag.Int("tcp-address", 4000, "tcp address")
	dir         = flag.String("dir", "", "directory")
)

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *tcp_address))
	if err != nil {
		panic(err)
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := master.Server(ln)
		if err != nil {
			fmt.Printf("master.Server server failed,err[%v]\n", err)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		runBenchMark()
	}()
	wg.Wait()
}

func runBenchMark() {
	s := 1024 * 1024
	for i := 0; i < s; i++ {

	}
}

func init() {
	master.RegisterGetter(&MemoryGetter{})
}

type MemoryGetter struct {
}

func (MemoryGetter) Get(k string) ([]byte, error) {
	length := crc32.ChecksumIEEE([]byte(k))
	length = length % (1024)
	s := time.Now().String() + k
	for uint32(len(s)) < length { // almost 1024 bytes
		s += s
	}
	return []byte(s), nil
}
