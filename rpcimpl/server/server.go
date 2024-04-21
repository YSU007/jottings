package server

import (
	"tiny_rpc/handler"
	"tiny_rpc/log"
	"tiny_rpc/module"
	"tiny_rpc/msg"
	"tiny_rpc/net"
	"tiny_rpc/router"
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
	msg.SetSerializer(msg.SerializerPB)
	handler.Init(router.ReflectRouterType)

	if err := module.MgrIns().Reg(q); err != nil {
		log.Error("reg module %v", err)
		return
	}
	if err := module.MgrIns().Start(); err != nil {
		log.Error("module mgr start %v", err)
		return
	}
	s.sm.Start()
}

func (s *Server) Close() {
	s.sm.Stop()

	if err := module.MgrIns().Stop(); err != nil {
		log.Error("server module mgr stop err %v", err)
		return
	}
}
