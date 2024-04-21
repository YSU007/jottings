package module

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"tiny_rpc/log"
	"tiny_rpc/util"
)

var m = newMgr()

func MgrIns() *Mgr {
	return m
}

func SyncWork(mo, me string, arg, reply interface{}) error {
	return MgrIns().work(mo, me, arg, reply)
}

func NotifyWork(mo, me string, arg interface{}) error {
	return MgrIns().work(mo, me, arg, nil)
}

// Que ----------------------------------------------------------------------------------------------------
type (
	Que  []Node
	Node struct {
		M     M
		Size  int
		MName string
	}
)

// Mgr ----------------------------------------------------------------------------------------------------
type Mgr struct {
	wg         *util.WGWrapper
	moduleMap  map[string]*Base
	moduleList []string
}

func newMgr() *Mgr {
	return &Mgr{
		wg:         new(util.WGWrapper),
		moduleMap:  make(map[string]*Base, 128),
		moduleList: make([]string, 0, 128),
	}
}

func (r *Mgr) Reg(q Que) error {
	for _, n := range q {
		t := reflect.TypeOf(n.M)
		if t.Kind() != reflect.Ptr {
			name := t.Name()
			if n.MName != "" {
				name = n.MName
			}
			return fmt.Errorf("reg module %s fail, module not ptr", name)
		}

		name := t.Elem().Name()
		if n.MName != "" {
			name = n.MName
		}
		if r.moduleMap[name] != nil {
			return fmt.Errorf("reg module %s fail, module duplicate", name)
		}

		module := NewModule(n.M, name, n.Size, r.wg)
		r.moduleMap[name] = module
		r.moduleList = append(r.moduleList, name)
		if err := module.register(); err != nil {
			return fmt.Errorf("reg module err %v", err)
		}
		log.Info("reg module %s chan size %d", name, n.Size)
	}
	return nil
}

func (r *Mgr) Start() error {
	for i := 0; i < len(r.moduleList); i++ {
		name := r.moduleList[i]
		if r.moduleMap[name] == nil {
			return fmt.Errorf("module mgr start %s err", name)
		}
		r.moduleMap[name].Start()
	}
	log.Info("modules start, count %d", len(r.moduleMap))
	return nil
}

func (r *Mgr) Stop() error {
	for i := len(r.moduleList) - 1; i >= 0; i-- {
		name := r.moduleList[i]
		if r.moduleMap[name] == nil {
			log.Error("module mgr stop %s err", name)
			continue
		}
		if err := r.moduleMap[name].Stop(); err != nil {
			log.Error("module %s stop err %v", name, err)
			continue
		}
	}
	log.Info("modules stop, count %d", len(r.moduleMap))
	return nil
}

func (r *Mgr) work(mo, me string, arg, reply interface{}) error {
	if r.moduleMap[mo] == nil {
		return fmt.Errorf("module %s nil", mo)
	}

	if reply == nil {
		r.moduleMap[mo].receive(&Work{
			Method: me,
			Arg:    arg,
		})
		return nil
	}

	w := &Work{
		Method:  me,
		Arg:     arg,
		Reply:   reply,
		RetChan: make(chan struct{}),
	}
	r.moduleMap[mo].receive(w)
	select {
	case <-w.RetChan:
		if w.Err != nil {
			return w.Err
		}
	}
	return nil
}

// M ----------------------------------------------------------------------------------------------------
type M interface {
	Load() error
	Save() error
}

// Base ----------------------------------------------------------------------------------------------------
type Base struct {
	name      string
	workChan  chan *Work
	closeChan chan interface{}
	wg        *util.WGWrapper
	rel       M
	typ       reflect.Type
	rcvr      reflect.Value
	method    map[string]*methodType
}

type Work struct {
	Method  string
	Arg     interface{}
	Reply   interface{}
	Err     error
	RetChan chan struct{}
}

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

func NewModule(m M, name string, chanSize int, wg *util.WGWrapper) *Base {
	return &Base{
		name:      name,
		workChan:  make(chan *Work, chanSize),
		closeChan: make(chan interface{}),
		wg:        wg,
		rel:       m,
		typ:       reflect.TypeOf(m),
		rcvr:      reflect.ValueOf(m),
	}
}

func (r *Base) Start() {
	defer util.InfoPanic("module %v panic", r.name)
	log.Info("module %v start", r.name)

	r.wg.Wrap(func() {
		for {
			select {
			case w := <-r.workChan:
				func(w *Work) {
					defer util.InfoPanic("module %v panic, work %s param %+v", r.name, w.Method, w.Arg)
					r.dealWork(w)
				}(w)
			case <-r.closeChan:
				return
			}
		}
	})
}

func (r *Base) Stop() error {
	close(r.closeChan)
	return nil
}

func (r *Base) receive(w *Work) {
	r.workChan <- w
}

func (r *Base) register() error {
	r.method = make(map[string]*methodType)
	for m := 0; m < r.typ.NumMethod(); m++ {
		method := r.typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}

		// notify method needs three ins: receiver, *args.
		if mtype.NumIn() == 2 {
			argType := mtype.In(1)

			if argType.Kind() != reflect.Ptr {
				return fmt.Errorf("arg not ptr")
			}
			r.method[mname] = &methodType{method: method, ArgType: argType}
			log.Info("module %s reg notify method %s", r.name, mname)
		}

		// sync method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() == 3 {
			argType := mtype.In(1)
			replyType := mtype.In(2)

			if argType.Kind() != reflect.Ptr || replyType.Kind() != reflect.Ptr {
				return fmt.Errorf("arg or reply not ptr")
			}
			// Method needs one out.
			if mtype.NumOut() != 1 {
				log.Info("method %s has wrong number of outs: %d", mname, mtype.NumOut())
				continue
			}
			r.method[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
			log.Info("module %s reg sync method %s", r.name, mname)
		}
	}
	return nil
}

func (r *Base) dealWork(w *Work) {
	mtype := r.method[w.Method]
	if mtype == nil {
		var err = fmt.Errorf("module %s method %s not find", r.name, w.Method)
		if w.RetChan != nil {
			w.Err = err
			w.RetChan <- struct{}{}
		}
		return
	}

	if w.Reply == nil || w.RetChan == nil {
		log.Debug("module %s receive notify %+v", r.name, w.Arg)
		r.callNotify(mtype, reflect.ValueOf(w.Arg))
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("module %s receive sync call arg %+v ", r.name, w.Arg))
	err := r.callSync(mtype, reflect.ValueOf(w.Arg), reflect.ValueOf(w.Reply))
	builder.WriteString(fmt.Sprintf("reply %+v err %v", w.Reply, err))
	log.Debug(builder.String())
	if w.RetChan != nil {
		if err != nil {
			w.Err = err
		}
		w.RetChan <- struct{}{}
	}
	return
}

func (r *Base) callSync(mtype *methodType, argv, replyv reflect.Value) (err error) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[module internal error]: %v, method: %s, argv: %+v, stack: %s",
				err, mtype.method.Name, argv.Interface(), buf)
			log.Error("%v", err)
		}
	}()

	function := mtype.method.Func
	returnValues := function.Call([]reflect.Value{r.rcvr, argv, replyv})
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}
	return nil
}

func (r *Base) callNotify(mtype *methodType, argv reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[module internal error]: %v, method: %s, argv: %+v, stack: %s",
				err, mtype.method.Name, argv.Interface(), buf)
			log.Error("%v", err)
		}
	}()

	function := mtype.method.Func
	function.Call([]reflect.Value{r.rcvr, argv})
}
