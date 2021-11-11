package mapping_handles

import (
	"Jottings/tiny_rpc/handler/mapping_handles/hello"
	"Jottings/tiny_rpc/proto"
	"Jottings/tiny_rpc/router"
)

// RegHandles ----------------------------------------------------------------------------------------------------
func RegHandles() {
	router.RegHandle(proto.Hello, hello.HelloProto{})
}
