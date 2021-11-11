package main

import (
	"time"

	"Jottings/tiny_rpc/client"
	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/msg"
	"Jottings/tiny_rpc/proto"
)

const (
	network = "tcp"
	address = "localhost:8972"
)

func main() {
	log.SetDefLog(log.FmtLog{})
	log.Info("Start..")
	defer log.Info("Stop..")
	msg.SetSerializer(msg.SerializerJson)
	var cli = client.NewClient(network, address)
	defer cli.Close()
	var req = &proto.HelloReq{
		HelloMsg: "hello",
	}
	var rsp = &proto.HelloRsp{}

	var ticker = time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			code, err := cli.Call(proto.Hello, req, rsp)
			if err != nil {
				return
			}
			log.Info("client call code %v %v rsp %v", code, err, rsp)
		}
	}
}
