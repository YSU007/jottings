package hello

import (
	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/model"
	"Jottings/tiny_rpc/proto"
)

func HelloHandle(a *model.PlayerAccount, req *proto.HelloReq, rsp *proto.HelloRsp) (code uint32) {
	log.Info("account %v receive %v", a.AccountId, req.HelloMsg)
	rsp.ReplyMsg = "hello~"
	return
}
