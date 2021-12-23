package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"rpcimpl/rpcserver/rpcproto"
)

// server is used to implement rpc methods.
type server struct {
	rpcproto.UnimplementedGameXServer
}

//RunRPCService simply run a rpc server.
func RunRPCService(netType string, lisPort string) {
	var lis, err = net.Listen(netType, lisPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var s = grpc.NewServer()
	rpcproto.RegisterGameXServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
