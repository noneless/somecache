package main

import (
	"flag"
	"sync"

	"github.com/756445638/somecache/slave"
)

var (
	master = flag.String("master-tcp-address", "127.0.0.1:4000", "master addr")
	worker = flag.Int("worker", 1, "worker count")
	wg     = sync.WaitGroup{}
)

func main() {
	flag.Parse()
	if *worker <= 0 {
		*worker = 1
	}
	if *worker > 10 {
		*worker = 10
	}
	wg.Add(*worker)
	for i := 0; i < *worker; i++ {
		go func() {
			slave.Connection2Master(*master)
			wg.Done()
		}()

	}
	wg.Wait()
}
