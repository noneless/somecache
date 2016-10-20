package grpc

import (
	"net"

	"github.com/75644563/somecache/grpc/connector"
	"google.golang.org/grpc"
)

type Grpc struct {
}

var (
	grpc Grpc
	ser  server
)

func RunMaster(ln net.Listener) error {
	return grpc.RunMaster(ln)
}
func (g *Grpc) RunMaster(net.Listener) error {
	s := grpc.NewServer()
	return nil
}

//lis, err := net.Listen("tcp", port)
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	s := grpc.NewServer()
//	pb.RegisterGreeterServer(s, &server{})
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
