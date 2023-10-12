package fsm

// Event 是在回调中作为引用传递的信息。
type Event struct {
	// FSM 是对当前FSM的引用。
	FSM *FSM

	// Event 是事件名称。
	Event string

	// Src 是转换前的状态。
	Src string

	// Dst 是转换后的状态。
	Dst string

	// Err 是可以从回调返回的可选错误。
	Err error

	// Args 是传递给回调的可选参数列表。
	Args []interface{}

	// canceled 是一个内部标志，如果转换被取消则设置。
	canceled bool

	// async 是一个内部标志设置，如果转换应该是异步的。
	async bool
}

// Cancel 可以在 before_<EVENT> 或 left_<STATE> 中调用以在当前转换发生之前取消当前转换。
// 它需要一个可选的错误，如果之前设置过，它将覆盖 e.Err。
func (e *Event) Cancel(err ...error) {
	e.canceled = true

	if len(err) > 0 {
		e.Err = err[0]
	}
}

// Async 可以在leave_<STATE>中调用Async来进行异步状态转换。
//
// 当前状态转换将保留在旧状态，直到最终状态转换
// 调用 Transition。这将完成过渡，并可能
// 调用其他回调。
func (e *Event) Async() {
	e.async = true
}
