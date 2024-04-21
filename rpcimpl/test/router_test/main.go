package main

import (
	"encoding/json"

	"tiny_rpc/handler"
	"tiny_rpc/log"
	"tiny_rpc/model"
	"tiny_rpc/msg"
	"tiny_rpc/proto"
	"tiny_rpc/router"
)

func main() {
	MappingRouter()
}

func MappingRouter() {
	log.Init(log.DefLog)

	handler.Init(router.MappingRouterType)

	var req = &proto.HelloReq{
		HelloMsg: "hello",
	}
	var baseReq = &msg.RequestBase{}
	var reqData, _ = json.Marshal(req)
	baseReq.FillIn(proto.Hello, reqData)

	var baseRsp = &msg.ResponseBase{}

	var a = &model.PlayerAccount{
		AccountId: "zhang",
	}
	router.HandleServe(a, baseReq, baseRsp)

	var rsp = &proto.HelloRsp{}
	var _ = json.Unmarshal(baseRsp.GetData(), rsp)
	log.Info("receive code %v data %v", baseRsp.GetCode(), *rsp)
}
