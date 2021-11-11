package router

import (
	"Jottings/tiny_rpc/msg"
)

type Type uint8

const (
	_ Type = iota
	MappingRouterType
	ReflectRouterType
)

// ContextInterface HandleInterface RouterInterface ----------------------------------------------------------------------------------------------------
type ContextInterface interface {
}

type HandleInterface interface {
	Serve(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg)
}

type Interface interface {
	RegHandle(mode uint32, handleInterface HandleInterface)
	HandleServe(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg)
}

// Instance ----------------------------------------------------------------------------------------------------
var Instance Interface

func SetRouterInstance(ins Interface) {
	Instance = ins
}

func RegHandle(mode uint32, handleInterface HandleInterface) {
	Instance.RegHandle(mode, handleInterface)
}

func HandleServe(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg) {
	Instance.HandleServe(ctx, req, rsp)
}
