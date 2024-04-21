package a

import (
	"tiny_rpc/log"
	"tiny_rpc/module"
	"tiny_rpc/proto"
)

type MA struct {
}

func (r *MA) Load() error {
	panic("implement me")
}

func (r *MA) Save() error {
	panic("implement me")
}

func (r *MA) Hello(arg *proto.HelloArg, replay *proto.HelloReplay) error {
	log.Info("hello %v", arg.Msg)
	replay.Msg = "hello~"
	err := module.NotifyWork(module.MB, module.MB_Hello, &proto.HelloArg{Msg: "mb"})
	if err != nil {
		return err
	}
	return nil
}

func (r *MA) Func1(arg *proto.HelloArg, replay *proto.HelloReplay) error {
	log.Info("MA Func1 %s", arg.Msg)
	replay.Msg = "Func1 hi"
	return nil
}
