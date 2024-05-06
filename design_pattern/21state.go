package main

import "fmt"

// State 定义状态接口
type State interface {
	DoAction(context *Context21)
}

// Context21 是一个拥有某种状态的类
type Context21 struct {
	State State
}

// 设置状态
func (c *Context21) SetState(state State) {
	c.State = state
}

// 执行状态行为
func (c *Context21) DoAction() {
	c.State.DoAction(c)
}

// StartState 表示开始状态
type StartState struct{}

func (s *StartState) DoAction(context *Context21) {
	fmt.Println("Player is in start state")
	context.SetState(s)
}

// StopState 表示停止状态
type StopState struct{}

func (s *StopState) DoAction(context *Context21) {
	fmt.Println("Player is in stop state")
	context.SetState(s)
}

func main() {
	context := new(Context21)

	startState := new(StartState)
	startState.DoAction(context)

	fmt.Println(context.State)

	stopState := new(StopState)
	stopState.DoAction(context)

	fmt.Println(context.State)
}
