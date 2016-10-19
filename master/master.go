package main

import (
	"flag"
)

var (
	master = flag.String("master", "", "master addr")
	tcp    = flag.String("tcp-address", "", "tcp address")
)

func main() {
	flag.Parse()
	if *master == "" {
		panic("master is nil string")
	}
	if *tcp == "" {
		panic("tcp is nil string")
	}
}

type Getter interface {
	Get(string) []byte
}
