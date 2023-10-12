package fsm

import (
	"strings"
	"sync"
)

// transitioner 是FSM转换函数的接口。
type transitioner interface {
	transition(*FSM) error
}

// FSM 是保存当前状态的状态机。
//
// 它必须使用 NewFSM 创建才能正常运行。
type FSM struct {
	// current is the state that the FSM is currently in.
	current string

	// transitions maps events and source states to destination states.
	transitions map[eKey]string

	// callbacks maps events and targets to callback functions.
	callbacks map[cKey]Callback

	// transition is the internal transition functions used either directly
	// or when Transition is called in an asynchronous state transition.
	transition func()
	// transitionerObj calls the FSM's transition() function.
	transitionerObj transitioner

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex
	// metadata can be used to store and load data that maybe used across events
	// use methods SetMetadata() and Metadata() to store and load data
	metadata map[string]interface{}

	metadataMu sync.RWMutex
}

// EventDesc 表示初始化 FSM 时的事件。
//
// 该事件可以有一个或多个对执行有效的源状态过渡。
// 如果 FSM 处于源状态之一，它将最终处于指定的目标状态，调用所有已定义的回调。
type EventDesc struct {
	// Name is the event name used when calling for a transition.
	Name string

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []string

	// Dst is the destination state that the FSM will be in if the transition
	// succeeds.
	Dst string
}

// Callback 是回调应该使用的函数类型。事件是当前的回调发生时的事件信息。
type Callback func(*Event)

// Events 是 NewFSM 中定义转换映射的简写。
type Events []EventDesc

// Callbacks 是 NewFSM 中定义回调的简写。
type Callbacks map[string]Callback

// NewFSM 根据事件和回调构造 FSM。
//
// 事件和转换被指定为事件结构的切片指定为事件。
// 每个事件都映射到一个或多个内部从Event.Src 转换到Event.Dst。
//
// 回调被添加为指定为解析键的回调的映射作为回调事件如下，并以相同的顺序调用：
//
// 1. before_<EVENT> - 在名为 <EVENT> 的事件之前调用
//
// 2. before_event - 在所有事件之前调用
//
// 3.leave_<OLD_STATE> - 在离开<OLD_STATE>之前调用
//
// 4.leave_state - 在离开所有状态之前调用
//
// 5. Enter_<NEW_STATE> - 输入 <NEW_STATE> 后调用
//
// 6. Enter_state - 进入所有状态后调用
//
// 7. after_<EVENT> - 在名为 <EVENT> 的事件之后调用
//
// 8. after_event - 在所有事件之后调用
//
// 最常用的回调还有两个简短版本。
// 它们只是事件或状态的名称：
//
// 1. <NEW_STATE> - 进入<NEW_STATE>后调用
//
// 2. <EVENT> - 在名为 <EVENT> 的事件之后调用
//
// 如果同时指定了简写版本和完整版本，则未定义哪个版本的回调将最终出现在内部映射中。
// 这个到期了Go Map的伪随机性。不检查多个Key当前执行时。
func NewFSM(initial string, events []EventDesc, callbacks map[string]Callback) *FSM {
	f := &FSM{
		transitionerObj: &transitionerStruct{},
		current:         initial,
		transitions:     make(map[eKey]string),
		callbacks:       make(map[cKey]Callback),
		metadata:        make(map[string]interface{}),
	}

	// Build transition map and store sets of all events and states.
	allEvents := make(map[string]bool)
	allStates := make(map[string]bool)
	for _, e := range events {
		for _, src := range e.Src {
			f.transitions[eKey{e.Name, src}] = e.Dst
			allStates[src] = true
			allStates[e.Dst] = true
		}
		allEvents[e.Name] = true
	}

	// Map all callbacks to events/states.
	for name, fn := range callbacks {
		var target string
		var callbackType int

		switch {
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "event" {
				target = ""
				callbackType = callbackBeforeEvent
			} else if _, ok := allEvents[target]; ok {
				callbackType = callbackBeforeEvent
			}
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "state" {
				target = ""
				callbackType = callbackLeaveState
			} else if _, ok := allStates[target]; ok {
				callbackType = callbackLeaveState
			}
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "state" {
				target = ""
				callbackType = callbackEnterState
			} else if _, ok := allStates[target]; ok {
				callbackType = callbackEnterState
			}
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "event" {
				target = ""
				callbackType = callbackAfterEvent
			} else if _, ok := allEvents[target]; ok {
				callbackType = callbackAfterEvent
			}
		default:
			target = name
			if _, ok := allStates[target]; ok {
				callbackType = callbackEnterState
			} else if _, ok := allEvents[target]; ok {
				callbackType = callbackAfterEvent
			}
		}

		if callbackType != callbackNone {
			f.callbacks[cKey{target, callbackType}] = fn
		}
	}

	return f
}

// Current returns the current state of the FSM.
func (f *FSM) Current() string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.current
}

// Is returns true if state is the current state.
func (f *FSM) Is(state string) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return state == f.current
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *FSM) SetState(state string) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = state
}

// Can returns true if event can occur in the current state.
func (f *FSM) Can(event string) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	_, ok := f.transitions[eKey{event, f.current}]
	return ok && (f.transition == nil)
}

// AvailableTransitions returns a list of transitions available in the
// current state.
func (f *FSM) AvailableTransitions() []string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	var transitions []string
	for key := range f.transitions {
		if key.src == f.current {
			transitions = append(transitions, key.event)
		}
	}
	return transitions
}

// Cannot returns true if event can not occur in the current state.
// It is a convenience method to help code read nicely.
func (f *FSM) Cannot(event string) bool {
	return !f.Can(event)
}

// Metadata returns the value stored in metadata
func (f *FSM) Metadata(key string) (interface{}, bool) {
	f.metadataMu.RLock()
	defer f.metadataMu.RUnlock()
	dataElement, ok := f.metadata[key]
	return dataElement, ok
}

// SetMetadata stores the dataValue in metadata indexing it with key
func (f *FSM) SetMetadata(key string, dataValue interface{}) {
	f.metadataMu.Lock()
	defer f.metadataMu.Unlock()
	f.metadata[key] = dataValue
}

// Event 使用指定事件启动状态转换。
//
// 该调用采用可变数量的参数，这些参数将传递给回调（如果已定义）。
//
// 如果状态更改正常或出现以下错误之一，它将返回 nil：
//
// - 事件 X 不合适，因为之前的转换未完成
//
// - 事件 X 在当前状态 Y 中不合适
//
// - 事件 X 不存在
//
// - 状态转换时的内部错误
//
// 在这种情况下，最后一个错误不应该发生，并且是内部错误的迹象。
func (f *FSM) Event(event string, args ...interface{}) error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()

	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	if f.transition != nil {
		return InTransitionError{event}
	}

	dst, ok := f.transitions[eKey{event, f.current}]
	if !ok {
		for ekey := range f.transitions {
			if ekey.event == event {
				return InvalidEventError{event, f.current}
			}
		}
		return UnknownEventError{event}
	}

	e := &Event{f, event, f.current, dst, nil, args, false, false}

	err := f.beforeEventCallbacks(e)
	if err != nil {
		return err
	}

	if f.current == dst {
		f.afterEventCallbacks(e)
		return NoTransitionError{e.Err}
	}

	// Setup the transition, call it later.
	f.transition = func() {
		f.stateMu.Lock()
		f.current = dst
		f.stateMu.Unlock()

		f.enterStateCallbacks(e)
		f.afterEventCallbacks(e)
	}

	if err = f.leaveStateCallbacks(e); err != nil {
		if _, ok := err.(CanceledError); ok {
			f.transition = nil
		}
		return err
	}

	// Perform the rest of the transition, if not asynchronous.
	f.stateMu.RUnlock()
	defer f.stateMu.RLock()
	err = f.doTransition()
	if err != nil {
		return InternalError{}
	}

	return e.Err
}

// Transition wraps transitioner.transition.
func (f *FSM) Transition() error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()
	return f.doTransition()
}

// doTransition wraps transitioner.transition.
func (f *FSM) doTransition() error {
	return f.transitionerObj.transition(f)
}

// transitionerStruct 是 transitioner 接口的默认实现。可以交换其他实现来进行测试。
type transitionerStruct struct{}

// transition 完成异步状态改变。
//
// leave_<STATE> 的回调之前必须在其事件上调用 Async 才能启动异步状态转换。
func (t transitionerStruct) transition(f *FSM) error {
	if f.transition == nil {
		return NotInTransitionError{}
	}
	f.transition()
	f.transition = nil
	return nil
}

// beforeEventCallbacks calls the before_ callbacks, first the named then the
// general version.
func (f *FSM) beforeEventCallbacks(e *Event) error {
	if fn, ok := f.callbacks[cKey{e.Event, callbackBeforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"", callbackBeforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the
// general version.
func (f *FSM) leaveStateCallbacks(e *Event) error {
	if fn, ok := f.callbacks[cKey{f.current, callbackLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"", callbackLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}
	return nil
}

// enterStateCallbacks calls the enter_ callbacks, first the named then the
// general version.
func (f *FSM) enterStateCallbacks(e *Event) {
	if fn, ok := f.callbacks[cKey{f.current, callbackEnterState}]; ok {
		fn(e)
	}
	if fn, ok := f.callbacks[cKey{"", callbackEnterState}]; ok {
		fn(e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the
// general version.
func (f *FSM) afterEventCallbacks(e *Event) {
	if fn, ok := f.callbacks[cKey{e.Event, callbackAfterEvent}]; ok {
		fn(e)
	}
	if fn, ok := f.callbacks[cKey{"", callbackAfterEvent}]; ok {
		fn(e)
	}
}

const (
	callbackNone int = iota
	callbackBeforeEvent
	callbackLeaveState
	callbackEnterState
	callbackAfterEvent
)

// cKey is a struct key used for keeping the callbacks mapped to a target.
type cKey struct {
	// target is either the name of a state or an event depending on which
	// callback type the key refers to. It can also be "" for a non-targeted
	// callback like before_event.
	target string

	// callbackType is the situation when the callback will be run.
	callbackType int
}

// eKey is a struct key used for storing the transition map.
type eKey struct {
	// event is the name of the event that the keys refers to.
	event string

	// src is the source from where the event can transition.
	src string
}
