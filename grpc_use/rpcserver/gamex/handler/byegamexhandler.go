package handler

import (
	"context"
	"fmt"

	"rpcimpl/rpcserver/rpcproto"
)

func ByeGameXHandler(ctx context.Context, in *rpcproto.ByeReq, out *rpcproto.ByeRsp) error {
	fmt.Println(in.Name, "say bye")
	out.Message = "bye " + in.Name
	return nil
}
