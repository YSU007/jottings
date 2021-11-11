package server

import (
	"Jottings/tiny_rpc/handler"
	"Jottings/tiny_rpc/net"
	"Jottings/tiny_rpc/router"
)

type Server struct {
	sm *net.SessionMgr
}

func NewServer(network, address string) *Server {
	return &Server{
		sm: net.NewSessionMgr(network, address),
	}
}

func (s *Server) Serve() {
	handler.Init(router.MappingRouterType)
	s.sm.Start()
}

func (s *Server) Close() {
	s.sm.Stop()
}
