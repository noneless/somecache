package grpc

import (
	"net"

	pb "github.com/756445638/somecache/master/grpc/proto"
	"google.golang.org/grpc"
)

var (
	ser server
)

func (g *server) RunMaster(net.Listener) error {
	s := grpc.NewServer()
	pb.RegisterMasterServer(s, ser)
	return nil
}

func RunMaster(ln net.Listener) error {
	return ser.RunMaster(ln)
}
