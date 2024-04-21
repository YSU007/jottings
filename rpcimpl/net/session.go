package net

import (
	"fmt"
	"io"
	"net"
	"strings"
	"unsafe"

	"tiny_rpc/log"
	"tiny_rpc/model"
	"tiny_rpc/msg"
	"tiny_rpc/router"
	"tiny_rpc/util"
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
	//db load
	s.Account = &model.PlayerAccount{AccountId: "zhang"}

	//go receive
	s.wg.Wrap(func() {
		for {
			var modeMsg = new(msg.ModeBase)
			var err = modeMsg.Decode(s)
			if err != nil {
				if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
					log.Info("Session %v connect close.", s.ID)
					return
				}
				log.Error("Session %d modeMsg Decode err %v", s.ID, err)
				return
			}

			s.works <- modeMsg
		}
	})

	//work handle
	s.handle()
}

func (s *Session) stop() {
	close(s.works)

	err := s.Close()
	if err != nil {
		log.Error("Session close err %v %v", s.ID, err)
		return
	}
}

func (s *Session) handle() {
	var err error
	for work := range s.works {
		switch work.MsgType() {
		case msg.MTypeRpc:
			err = s.handleRPC((*msg.RequestBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		case msg.MTypeNotice:
			s.handleNotify((*msg.NotifyBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		case msg.MTypePush:
			s.handlePush((*msg.PushBase)(unsafe.Pointer(work.(*msg.ModeBase))))
		}
		if err != nil {
			log.Error("Session %d err %v", s.ID, err)
			s.stop()
		}
	}
}

func (s *Session) handleRPC(baseReq *msg.RequestBase) error {
	var baseRsp = new(msg.ResponseBase)

	// serve handle
	var err = router.HandleServe(s.Account, baseReq, baseRsp)
	if err != nil {
		return err
	}

	err = baseRsp.Encode(s)
	if err != nil {
		if err == io.EOF {
			log.Info("Session %v connect close.", s.ID)
			return nil
		}
		return fmt.Errorf("write err %v", err)
	}
	return nil
}

func (s *Session) handleNotify(baseNotify *msg.NotifyBase) {
}

func (s *Session) handlePush(basePush *msg.PushBase) {
}
