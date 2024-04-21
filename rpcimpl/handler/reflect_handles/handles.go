package reflect_handles

import (
	"tiny_rpc/handler/reflect_handles/hello"
	"tiny_rpc/proto"
	"tiny_rpc/router"
)

// RegHandlesFunc ----------------------------------------------------------------------------------------------------
func RegHandlesFunc() {
	router.RegHandle(proto.Hello, router.NewFuncHandle(hello.HelloHandle))
}
