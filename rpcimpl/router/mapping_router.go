package router

import (
	"fmt"

	"tiny_rpc/msg"
)

// MappingRouter ----------------------------------------------------------------------------------------------------
type MappingRouter map[uint32]HandleInterface

func (r *MappingRouter) RegHandle(mode uint32, handleInterface HandleInterface) {
	(*r)[mode] = handleInterface
}

func (r *MappingRouter) HandleServe(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg) error {
	var mode = req.GetMode()
	var f = (*r)[mode]
	if f == nil {
		return fmt.Errorf("mode %d not find", req.GetMode())
	}
	f.Serve(ctx, req, rsp)
	return nil
}
