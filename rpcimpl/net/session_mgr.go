package net

import (
	"net"
	"sync"
	"sync/atomic"

	"tiny_rpc/log"
	"tiny_rpc/util"
)

type SessionMgr struct {
	net.Listener
	wg             *util.WGWrapper
	sessionCounter SessionID
	l              sync.RWMutex
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

func (r *SessionMgr) Start() {
	for {
		conn, err := r.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Error("server accept err %v", err)
				continue
			}
			log.Info("server accept err %v", err)
			return
		}
		atomic.AddUint32((*uint32)(&(r.sessionCounter)), 1)
		var session = newSession(conn, r.sessionCounter, r.wg)
		r.add(session)
		r.wg.Wrap(session.start)
	}
}

func (r *SessionMgr) Stop() {
	err := r.Close()
	if err != nil {
		log.Error("server stop err %v", err)
	}

	var counter = func() int {
		r.l.RLock()
		defer r.l.RUnlock()

		var counter = len(r.sessions)
		for _, s := range r.sessions {
			s.stop()
		}
		return counter
	}()
	r.wg.Wait()
	log.Info("server stop,accept count %d,stop session %d.", r.sessionCounter, counter)
}

func (r *SessionMgr) add(s *Session) {
	r.l.Lock()
	defer r.l.Unlock()
	r.sessions[r.sessionCounter] = s
}
