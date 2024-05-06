package main

import "fmt"

// Memento 包含了要被恢复的对象的状态
type Memento struct {
	state string
}

// Originator 是拥有当前状态的对象
type Originator struct {
	state string
}

// 创建备忘录，保存当前状态
func (o *Originator) CreateMemento() *Memento {
	return &Memento{state: o.state}
}

// 从备忘录恢复状态
func (o *Originator) RestoreMemento(m *Memento) {
	o.state = m.state
}

// Caretaker 对备忘录进行管理，但不知道其内容
type Caretaker struct {
	memento *Memento
}

func main() {
	originator := &Originator{state: "On"}
	fmt.Println("初始状态:", originator.state)

	// 保存状态
	caretaker := &Caretaker{
		memento: originator.CreateMemento(),
	}

	// 改变状态
	originator.state = "Off"
	fmt.Println("新的状态:", originator.state)

	// 恢复状态
	originator.RestoreMemento(caretaker.memento)
	fmt.Println("恢复后的状态:", originator.state)
}
