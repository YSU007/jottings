package fsm

// InvalidEventError 当当前状态下无法调用事件时，由 FSM.Event() 返回。
type InvalidEventError struct {
	Event string
	State string
}

func (e InvalidEventError) Error() string {
	return "event " + e.Event + " inappropriate in current state " + e.State
}

// UnknownEventError 当事件未定义时，由 FSM.Event() 返回。
type UnknownEventError struct {
	Event string
}

func (e UnknownEventError) Error() string {
	return "event " + e.Event + " does not exist"
}

// InTransitionError 当异步转换已经在进行时，由 FSM.Event() 返回。
type InTransitionError struct {
	Event string
}

func (e InTransitionError) Error() string {
	return "event " + e.Event + " inappropriate because previous transition did not complete"
}

// NotInTransitionError 当异步转换未进行时，由 FSM.Transition() 返回。
type NotInTransitionError struct{}

func (e NotInTransitionError) Error() string {
	return "transition inappropriate because no state change in progress"
}

// NoTransitionError 当没有发生转换时，例如源状态和目标状态相同时，由 FSM.Event() 返回。
type NoTransitionError struct {
	Err error
}

func (e NoTransitionError) Error() string {
	if e.Err != nil {
		return "no transition with error: " + e.Err.Error()
	}
	return "no transition"
}

// CanceledError 当回调取消转换时，由 FSM.Event() 返回。
type CanceledError struct {
	Err error
}

func (e CanceledError) Error() string {
	if e.Err != nil {
		return "transition canceled with error: " + e.Err.Error()
	}
	return "transition canceled"
}

// AsyncError 当回调启动异步状态转换时，由 FSM.Event() 返回。
type AsyncError struct {
	Err error
}

func (e AsyncError) Error() string {
	if e.Err != nil {
		return "async started with error: " + e.Err.Error()
	}
	return "async started"
}

// InternalError 由 FSM.Event() 返回并且永远不应该发生。这可能是因为一个错误。
type InternalError struct{}

func (e InternalError) Error() string {
	return "internal error on state transition"
}
