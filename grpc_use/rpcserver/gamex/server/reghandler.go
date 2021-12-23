package server

import (
	"context"

	"rpcimpl/rpcserver/rpcproto"
	"rpcimpl/rpcserver/rpcserver/gamex/handler"
)

func (s *server) HelloGameX(ctx context.Context, in *rpcproto.HelloReq) (*rpcproto.HelloRsp, error) {
	var out = &rpcproto.HelloRsp{}
	return out, handler.HelloGameXHandler(ctx, in, out)
}

func (s *server) ByeGameX(ctx context.Context, in *rpcproto.ByeReq) (*rpcproto.ByeRsp, error) {
	var out = &rpcproto.ByeRsp{}
	return out, handler.ByeGameXHandler(ctx, in, out)
}
