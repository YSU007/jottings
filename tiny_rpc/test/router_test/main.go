package main

import (
	"encoding/json"

	"Jottings/tiny_rpc/handler"
	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/model"
	"Jottings/tiny_rpc/msg"
	"Jottings/tiny_rpc/proto"
	"Jottings/tiny_rpc/router"
)

func main() {
	MappingRouter()
}

func MappingRouter() {
	log.SetDefLog(log.FmtLog{})

	handler.Init(router.MappingRouterType)

	var req = &proto.HelloReq{
		HelloMsg: "hello",
	}
	var baseReq = &msg.RequestBase{}
	var reqData, _ = json.Marshal(req)
	baseReq.FillIn(handler.Hello, reqData)

	var baseRsp = &msg.ResponseBase{}

	var a = &model.PlayerAccount{
		AccountId: "zhang",
	}
	router.HandleServe(a, baseReq, baseRsp)

	var rsp = &proto.HelloRsp{}
	var _ = json.Unmarshal(baseRsp.GetData(), rsp)
	log.Info("receive code %v data %v", baseRsp.GetCode(), *rsp)
}
