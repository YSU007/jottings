package main

import "fmt"

// Visitor 接口定义了访问者的行为
type Visitor interface {
	VisitConcreteElementA(*ConcreteElementA)
	VisitConcreteElementB(*ConcreteElementB)
}

// Element 接口定义了元素的行为
type Element interface {
	Accept(Visitor)
}

// ConcreteElementA 是具体的元素类
type ConcreteElementA struct {
	value string
}

// Accept 方法接受访问者访问
func (e *ConcreteElementA) Accept(visitor Visitor) {
	visitor.VisitConcreteElementA(e)
}

// ConcreteElementB 是具体的元素类
type ConcreteElementB struct {
	value string
}

// Accept 方法接受访问者访问
func (e *ConcreteElementB) Accept(visitor Visitor) {
	visitor.VisitConcreteElementB(e)
}

// ConcreteVisitor 是具体的访问者类
type ConcreteVisitor struct{}

// VisitConcreteElementA 是 ConcreteVisitor 访问 ConcreteElementA 的方法
func (v *ConcreteVisitor) VisitConcreteElementA(e *ConcreteElementA) {
	fmt.Println("Visiting ConcreteElementA:", e.value)
}

// VisitConcreteElementB 是 ConcreteVisitor 访问 ConcreteElementB 的方法
func (v *ConcreteVisitor) VisitConcreteElementB(e *ConcreteElementB) {
	fmt.Println("Visiting ConcreteElementB:", e.value)
}

func main() {
	elementA := &ConcreteElementA{value: "A"}
	elementB := &ConcreteElementB{value: "B"}
	visitor := &ConcreteVisitor{}

	elementA.Accept(visitor)
	elementB.Accept(visitor)
}
