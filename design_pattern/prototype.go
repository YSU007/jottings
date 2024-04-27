package main

import "fmt"

// 定义接口
type Prototype interface {
	Clone() Prototype
}

// 定义结构体
type ConcretePrototype struct {
	name string
}

func (p *ConcretePrototype) Clone() Prototype {
	return &ConcretePrototype{name: p.name + "_clone"}
}

func main() {
	// 创建原型对象
	prototype := &ConcretePrototype{name: "Original"}

	// 克隆对象
	clone1 := prototype.Clone()
	clone2 := prototype.Clone()

	fmt.Println("Original:", prototype.name)
	fmt.Println("Clone 1:", clone1.(*ConcretePrototype).name)
	fmt.Println("Clone 2:", clone2.(*ConcretePrototype).name)
}
