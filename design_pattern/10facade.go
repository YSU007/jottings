package main

import "fmt"

// 子系统1
type SubSystemOne struct{}

func (s *SubSystemOne) OperationOne() {
	fmt.Println("SubSystemOne method")
}

// 子系统2
type SubSystemTwo struct{}

func (s *SubSystemTwo) OperationTwo() {
	fmt.Println("SubSystemTwo method")
}

// 子系统3
type SubSystemThree struct{}

func (s *SubSystemThree) OperationThree() {
	fmt.Println("SubSystemThree method")
}

// 外观类，用户只需要与这个类交互
type Facade struct {
	one   *SubSystemOne
	two   *SubSystemTwo
	three *SubSystemThree
}

func NewFacade() *Facade {
	return &Facade{
		one:   &SubSystemOne{},
		two:   &SubSystemTwo{},
		three: &SubSystemThree{},
	}
}

// 外观类的方法，它将客户端的请求代理给适当的子系统对象
func (f *Facade) OperationWrapper() {
	fmt.Println("OperationWrapper method called")
	f.one.OperationOne()
	f.two.OperationTwo()
	f.three.OperationThree()
}

func main() {
	facade := NewFacade()
	facade.OperationWrapper()
}
