package main

import "fmt"

// Mediator defines the interface for mediators
type Mediator interface {
	SendMessage(msg string, user User)
	ReceiveMessage() string
}

// ChatRoom is a concrete mediator
type ChatRoom struct {
	Message string
}

func (c *ChatRoom) SendMessage(msg string, user User) {
	c.Message = fmt.Sprintf("%s sent a message: %s\n", user.name, msg)
}

func (c *ChatRoom) ReceiveMessage() string {
	return c.Message
}

// User represents a user
type User struct {
	name     string
	mediator Mediator
}

func (u *User) SetMediator(mediator Mediator) {
	u.mediator = mediator
}

func (u *User) SendMessage(msg string) {
	u.mediator.SendMessage(msg, *u)
}

func (u *User) ReceiveMessage() string {
	return u.mediator.ReceiveMessage()
}

func main() {
	chatRoom := &ChatRoom{}

	user1 := &User{name: "User1"}
	user2 := &User{name: "User2"}

	user1.SetMediator(chatRoom)
	user2.SetMediator(chatRoom)

	user1.SendMessage("Hello, User2!")
	fmt.Println(user2.ReceiveMessage())

	user2.SendMessage("Hi, User1!")
	fmt.Println(user1.ReceiveMessage())
}
