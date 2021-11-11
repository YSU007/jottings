package handler

import (
	"Jottings/tiny_rpc/handler/mapping_handles"
	"Jottings/tiny_rpc/handler/reflect_handles"
	"Jottings/tiny_rpc/router"
)

func Init(t router.Type) {
	switch t {
	case router.MappingRouterType:
		{
			var DefRouter = make(router.MappingRouter)
			router.SetRouterInstance(&DefRouter)
			mapping_handles.RegHandles()
		}
	case router.ReflectRouterType:
		{
			var DefRouter = router.NewReflectRouter()
			router.SetRouterInstance(DefRouter)
			reflect_handles.RegHandlesFunc()
		}
	}
}
