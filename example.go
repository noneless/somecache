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
	s := 1024
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
	fmt.Println("################# get 1024 * 1024 objects,takes:", total.Seconds())

}

func init() {
	master.RegisterGetter(&MemoryGetter{})
}

type MemoryGetter struct {
}

func (MemoryGetter) Get(k string) ([]byte, error) {
	length := crc32.ChecksumIEEE([]byte(k))
	length = length % (4 * 1024 * 1024)
	s := time.Now().String() + k
	for uint32(len(s)) < length { // less than 4 * 1024
		s += s
	}
	return []byte(s), nil
}
