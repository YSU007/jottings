package router

import "Jottings/tiny_rpc/msg"

// MappingRouter ----------------------------------------------------------------------------------------------------
type MappingRouter map[uint32]HandleInterface

func (r *MappingRouter) RegHandle(mode uint32, handleInterface HandleInterface) {
	(*r)[mode] = handleInterface
}

func (r *MappingRouter) HandleServe(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg) {
	var mode = req.GetMode()
	(*r)[mode].Serve(ctx, req, rsp)
}
