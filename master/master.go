package main

import (
	"flag"
	"net"
	"sync"

	"github.com/756445638/somecache/master/grpc"
)

var (
	wg   = sync.WaitGroup{}
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
