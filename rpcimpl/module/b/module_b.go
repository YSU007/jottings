package b

import (
	"tiny_rpc/log"
	"tiny_rpc/module"
	"tiny_rpc/proto"
)

type MB struct {
}

func (r *MB) Load() error {
	panic("implement me")
}

func (r *MB) Save() error {
	panic("implement me")
}

func (r *MB) Hello(arg *proto.HelloArg) {
	log.Info("hello %s", arg.Msg)
	replay := &proto.HelloReplay{}
	module.SyncWork(module.MA, module.MA_Func1, &proto.HelloArg{Msg: "ma"}, replay)
	log.Info("MB %v", replay.Msg)
}
