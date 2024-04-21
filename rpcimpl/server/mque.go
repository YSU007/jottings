package server

import (
	"tiny_rpc/module"
	"tiny_rpc/module/a"
	"tiny_rpc/module/b"
)

var q = module.Que{
	{M: &a.MA{}, Size: module.ChanSizeDef, MName: module.MA},
	{M: &b.MB{}, Size: module.ChanSizeDef, MName: module.MB},
}
