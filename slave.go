package main

import (
	"flag"
	"sync"
	"time"

	"github.com/756445638/somecache/slave"
)

var (
	master = flag.String("master-tcp-address", "127.0.0.1:4000", "master addr")
	worker = flag.Int("worker", 1, "worker count")
	wg     = sync.WaitGroup{}
)

func main() {
	flag.Parse()
	if *worker <= 3 {
		*worker = 3
	}
	if *worker > 30 {
		*worker = 30
	}
	wg.Add(*worker)
	for i := 0; i < *worker; i++ {
		time.Sleep(time.Second)
		go func() {
			slave.Connection2Master(*master)
			wg.Done()
		}()
	}
	wg.Wait()
}
