package net

import (
	"net"
	"sync/atomic"

	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/util"
)

type SessionMgr struct {
	net.Listener
	wg             *util.WGWrapper
	sessionCounter SessionID
	sessions       map[SessionID]*Session
}

func NewSessionMgr(network, address string) *SessionMgr {
	listen, err := net.Listen(network, address)
	if err != nil {
		log.Error("NewSessionMgr err %v", err)
		return nil
	}
	return &SessionMgr{
		Listener: listen,
		wg:       new(util.WGWrapper),
		sessions: make(map[SessionID]*Session, 2^12),
	}
}

func (s *SessionMgr) Start() {
	for {
		conn, err := s.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Error("server accept err %v", err)
				continue
			}
			log.Info("server accept err %v", err)
			return
		}
		atomic.AddUint32((*uint32)(&(s.sessionCounter)), 1)
		var session = newSession(conn, s.sessionCounter, s.wg)
		s.sessions[s.sessionCounter] = session
		s.wg.Wrap(session.start)
	}
}

func (s *SessionMgr) Stop() {
	err := s.Close()
	if err != nil {
		log.Error("server stop err %v", err)
	}
	var counter = len(s.sessions)
	for _, s := range s.sessions {
		s.stop()
	}
	s.wg.Wait()
	log.Info("server stop,accept count %d,stop session %d.", s.sessionCounter, counter)
}
