package mapping_handles

import (
	"tiny_rpc/handler/mapping_handles/hello"
	"tiny_rpc/proto"
	"tiny_rpc/router"
)

// RegHandles ----------------------------------------------------------------------------------------------------
func RegHandles() {
	router.RegHandle(proto.Hello, hello.HelloProto{})
}
