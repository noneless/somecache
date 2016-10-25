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
	"flag"
	"sync"
	"time"

	"github.com/756445638/somecache/slave"
)

var (
	master    = flag.String("master-tcp-address", "127.0.0.1:4000", "master addr")
	worker    = flag.Int("worker", 4, "worker count")
	cachesize = flag.Int64("cachesize", 1024, "cachesize in MB")

	wg = sync.WaitGroup{}
)

func main() {
	flag.Parse()
	if *worker <= 1 {
		*worker = 1
	}
	if *worker > 5 {
		*worker = 5
	}
	wg.Add(*worker)
	for i := 0; i < *worker; i++ {
		time.Sleep(time.Second)
		go func() {
			defer wg.Done()
			slave.Connection2Master(*master, (*cachesize)<<20)
		}()
	}
	wg.Wait()
}
