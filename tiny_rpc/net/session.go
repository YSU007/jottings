package net

import (
	"io"
	"net"
	"unsafe"

	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/model"
	"Jottings/tiny_rpc/msg"
	"Jottings/tiny_rpc/router"
	"Jottings/tiny_rpc/util"
)

type SessionID uint32

type Session struct {
	net.Conn
	wg      *util.WGWrapper
	works   chan msg.ModeMsg
	ID      SessionID
	Account model.AccountI
}

func newSession(conn net.Conn, id SessionID, wg *util.WGWrapper) *Session {
	return &Session{
		Conn:  conn,
		ID:    id,
		wg:    wg,
		works: make(chan msg.ModeMsg, 2^10),
	}
}

func (s *Session) start() {
	//go receive
	s.wg.Wrap(func() {
		for {
			var modeMsg = new(msg.ModeBase)
			var err = modeMsg.Decode(s)
			if err != nil {
				if err == io.EOF {
					log.Info("Session %v connect close.", s.ID)
					return
				}
				log.Error("Session modeMsg Decode err %v", err)
				return
			}

			s.works <- modeMsg
		}
	})

	//work handle
	s.handle()
}

func (s *Session) stop() {
	err := s.Close()
	if err != nil {
		log.Error("Session close err %v %v", s.ID, err)
		return
	}
}

func (s *Session) handle() {
	for work := range s.works {
		switch work.MsgType() {
		case msg.MTypeRpc:
			s.handleRPC((*msg.RequestBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		case msg.MTypeNotice:
			s.handleNotify((*msg.NotifyBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		case msg.MTypePush:
			s.handlePush((*msg.PushBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		}
	}
}

func (s *Session) handleRPC(baseReq *msg.RequestBase) {
	var baseRsp = new(msg.ResponseBase)
	var a = &model.PlayerAccount{AccountId: "zhang"}

	// serve handle
	router.HandleServe(a, baseReq, baseRsp)

	var err = baseRsp.Encode(s)
	if err != nil {
		if err == io.EOF {
			log.Info("Session %v connect close.", s.ID)
			return
		}
		log.Error("Session Write err %v", err)
		return
	}
}

func (s *Session) handleNotify(baseNotify *msg.NotifyBase) {
}

func (s *Session) handlePush(basePush *msg.PushBase) {
}
