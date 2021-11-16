package router

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/model"
	"Jottings/tiny_rpc/msg"
)

// funcHandle ----------------------------------------------------------------------------------------------------
type funcHandle struct {
	funcV     reflect.Value
	ArgType   reflect.Type
	ReplyType reflect.Type
}

func NewFuncHandle(fn interface{}) *funcHandle {
	f, ok := fn.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(fn)
	}
	if f.Kind() != reflect.Func {
		log.Error("function must be func or bound method")
		return nil
	}

	t := f.Type()
	if t.NumIn() != 3 {
		log.Error("registerFunction: has wrong number of ins: %r", f.Type().String())
		return nil
	}
	if t.NumOut() != 1 {
		log.Error("registerFunction: has wrong number of outs: %r", f.Type().String())
		return nil
	}

	argType := t.In(1)
	replyType := t.In(2)

	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)

	return &funcHandle{funcV: f, ArgType: argType, ReplyType: replyType}
}

func (r *funcHandle) Serve(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg) {
	argi := reflectTypePools.Get(r.ArgType)
	defer reflectTypePools.Put(r.ArgType, argi)
	replyi := reflectTypePools.Get(r.ReplyType)
	defer reflectTypePools.Put(r.ReplyType, replyi)

	// req unmarshal
	if err := msg.Unmarshal(req.GetData(), argi); err != nil {
		log.Error("funcHandle Serve Unmarshal err %v", err)
		return
	}

	//todo strings.Builder{}
	var tempStr string
	argv := reflect.ValueOf(argi)
	replyv := reflect.ValueOf(replyi)
	var a, ok = ctx.(model.AccountI)
	if ok && r.ArgType.Kind() == reflect.Ptr {
		tempStr += fmt.Sprintf("account %v mode %v request %+v", a.ID(), req.GetMode(), argv.Elem())
	}
	var code = r.call(ctx, argv, replyv)
	if ok && r.ReplyType.Kind() == reflect.Ptr {
		tempStr += fmt.Sprintf(" code %v request %+v", code, replyv.Elem())
		log.Debug(tempStr)
	}

	// rsp marshal
	var data, err = msg.Marshal(replyi)
	if err != nil {
		log.Error("funcHandle Serve Marshal err %v", err)
		return
	}
	rsp.FillIn(code, data)
}

func (r *funcHandle) call(ctx ContextInterface, argv, replyv reflect.Value) uint32 {
	fh := r

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err := fmt.Errorf("[callForFunction error]: %v, function: %r, argv: %+v, stack: %r",
				r, runtime.FuncForPC(fh.funcV.Pointer()), argv.Interface(), buf)
			log.Error("callForFunction %v", err)
		}
	}()

	returnValues := fh.funcV.Call([]reflect.Value{reflect.ValueOf(ctx), argv, replyv})
	code := returnValues[0].Interface()
	if code != nil {
		return code.(uint32)
	}
	return 0
}

// ReflectRouter ----------------------------------------------------------------------------------------------------
type ReflectRouter struct {
	function map[uint32]HandleInterface // registered functions
}

func NewReflectRouter() *ReflectRouter {
	return &ReflectRouter{
		function: make(map[uint32]HandleInterface),
	}
}

func (r *ReflectRouter) RegHandle(mode uint32, handleInterface HandleInterface) {
	r.function[mode] = handleInterface
}

func (r *ReflectRouter) HandleServe(ctx ContextInterface, req msg.ModeMsg, rsp msg.CodeMsg) {
	f := r.function[req.GetMode()]
	if f != nil {
		f.Serve(ctx, req, rsp)
	}
}

// typePools ----------------------------------------------------------------------------------------------------
type typePools struct {
	pools map[reflect.Type]*sync.Pool
	New   func(t reflect.Type) interface{}
}

func (p *typePools) Init(t reflect.Type) {
	tp := &sync.Pool{}
	tp.New = func() interface{} {
		return p.New(t)
	}
	p.pools[t] = tp
}

func (p *typePools) Put(t reflect.Type, x interface{}) {
	if o, ok := x.(Reset); ok {
		o.Reset()
	}
	pool := p.pools[t]
	pool.Put(x)
}

func (p *typePools) Get(t reflect.Type) interface{} {
	pool := p.pools[t]
	return pool.Get()
}

var reflectTypePools = &typePools{
	pools: make(map[reflect.Type]*sync.Pool),
	New: func(t reflect.Type) interface{} {
		var argv reflect.Value

		if t.Kind() == reflect.Ptr {
			argv = reflect.New(t.Elem())
		} else {
			argv = reflect.New(t)
		}

		return argv.Interface()
	},
}
