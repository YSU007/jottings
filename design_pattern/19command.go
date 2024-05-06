package main

import "fmt"

// Command 是命令接口
type Command interface {
	Execute()
}

// Light 是接收者类
type Light struct {
	isOn bool
}

// TurnOnCommand 是打开灯的命令
type TurnOnCommand struct {
	light *Light
}

func (c *TurnOnCommand) Execute() {
	c.light.isOn = true
	fmt.Println("Light is on")
}

// TurnOffCommand 是关闭灯的命令
type TurnOffCommand struct {
	light *Light
}

func (c *TurnOffCommand) Execute() {
	c.light.isOn = false
	fmt.Println("Light is off")
}

// Switch 是调用者类
type Switch struct {
	command Command
}

func (s *Switch) Press() {
	s.command.Execute()
}

func main() {
	light := &Light{}

	turnOn := &TurnOnCommand{light: light}
	turnOff := &TurnOffCommand{light: light}

	switchOn := &Switch{command: turnOn}
	switchOff := &Switch{command: turnOff}

	switchOn.Press()
	switchOff.Press()
}
