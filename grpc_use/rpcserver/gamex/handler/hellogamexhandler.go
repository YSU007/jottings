package handler

import (
	"context"
	"fmt"

	"rpcimpl/rpcserver/rpcproto"
)

func HelloGameXHandler(ctx context.Context, in *rpcproto.HelloReq, out *rpcproto.HelloRsp) error {
	fmt.Println(in.Name, "say hello")
	out.Message = "hello " + in.Name
	return nil
}
