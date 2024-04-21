package hello

import (
	"tiny_rpc/log"
	"tiny_rpc/model"
	"tiny_rpc/module"
	"tiny_rpc/proto"
)

func HelloHandle(a *model.PlayerAccount, req *proto.HelloReq, rsp *proto.HelloRsp) (code uint32) {
	log.Info("account %v receive %v", a.AccountId, req.HelloMsg)

	replay := &proto.HelloReplay{}
	err := module.SyncWork(module.MA, module.MA_Hello, &proto.HelloArg{Msg: "ma"}, replay)
	if err != nil {
		log.Error("%v", err)
	}
	rsp.ReplyMsg = replay.Msg
	return
}
