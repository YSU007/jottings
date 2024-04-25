package main

import "fmt"

// Strategy 接口定义了一个策略或行为的协议
type Strategy interface {
	DoOperation(num1 int, num2 int) int
}

// OperationAdd 实现了 Strategy 接口的加法策略
type OperationAdd struct{}

func (o *OperationAdd) DoOperation(num1 int, num2 int) int {
	return num1 + num2
}

// OperationSubstract 实现了 Strategy 接口的减法策略
type OperationSubstract struct{}

func (o *OperationSubstract) DoOperation(num1 int, num2 int) int {
	return num1 - num2
}

// OperationMultiply 实现了 Strategy 接口的乘法策略
type OperationMultiply struct{}

func (o *OperationMultiply) DoOperation(num1 int, num2 int) int {
	return num1 * num2
}

// Context 是一个使用了某种策略的类
type Context struct {
	strategy Strategy
}

func NewContext(strategy Strategy) *Context {
	return &Context{
		strategy: strategy,
	}
}

func (c *Context) ExecuteStrategy(num1 int, num2 int) int {
	return c.strategy.DoOperation(num1, num2)
}

func main() {
	context := NewContext(new(OperationAdd))
	fmt.Println("10 + 5 = ", context.ExecuteStrategy(10, 5))

	context = NewContext(new(OperationSubstract))
	fmt.Println("10 - 5 = ", context.ExecuteStrategy(10, 5))

	context = NewContext(new(OperationMultiply))
	fmt.Println("10 * 5 = ", context.ExecuteStrategy(10, 5))
}
