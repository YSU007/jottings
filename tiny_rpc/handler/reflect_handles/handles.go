package reflect_handles

import (
	"Jottings/tiny_rpc/handler/reflect_handles/hello"
	"Jottings/tiny_rpc/proto"
	"Jottings/tiny_rpc/router"
)

// RegHandlesFunc ----------------------------------------------------------------------------------------------------
func RegHandlesFunc() {
	router.RegHandle(proto.Hello, router.NewFuncHandle(hello.HelloHandle))
}
