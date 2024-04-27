package main

import "fmt"

// Computer 是客户端期望的接口
type Computer interface {
	InsertIntoLightningPort()
}

// Mac 是 Mac 电脑的实现
type Mac struct{}

func (m *Mac) InsertIntoLightningPort() {
	fmt.Println("Lightning connector is plugged into Mac machine.")
}

// Windows 是 Windows 笔记本的实现
type Windows struct{}

func (w *Windows) insertIntoUSBPort() {
	fmt.Println("USB connector is plugged into Windows machine.")
}

// WindowsAdapter 是适配器，将 Lightning 接口转换为 USB 接口
type WindowsAdapter struct {
	windowMachine *Windows
}

func (w *WindowsAdapter) InsertIntoLightningPort() {
	fmt.Println("Adapter converts Lightning signal to USB.")
	w.windowMachine.insertIntoUSBPort()
}

// Client 是客户端代码
type Client struct{}

func (c *Client) InsertLightningConnectorIntoComputer(com Computer) {
	fmt.Println("Client inserts Lightning connector into computer.")
	com.InsertIntoLightningPort()
}

func main() {
	client := &Client{}
	mac := &Mac{}
	client.InsertLightningConnectorIntoComputer(mac)

	windowsMachine := &Windows{}
	windowsMachineAdapter := &WindowsAdapter{windowMachine: windowsMachine}
	client.InsertLightningConnectorIntoComputer(windowsMachineAdapter)
}
