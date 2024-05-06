package main

import "fmt"

// Flyweight 接口定义享元对象的方法
type Flyweight interface {
	Operation(extrinsicState string)
}

// ConcreteFlyweight 是具体的享元类，包含内部状态和外部状态
type ConcreteFlyweight struct {
	color string
}

// Operation 方法实现在 ConcreteFlyweight 中搜索关键词的功能
func (c *ConcreteFlyweight) Operation(extrinsicState string) {
	fmt.Printf("Searching for keyword %s in color %s\n", extrinsicState, c.color)
}

func main() {
	redFlyweight := &ConcreteFlyweight{color: "red"}
	blueFlyweight := &ConcreteFlyweight{color: "blue"}

	redFlyweight.Operation("apple")
	blueFlyweight.Operation("sky")
}
