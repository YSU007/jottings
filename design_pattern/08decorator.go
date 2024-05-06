package main

import "fmt"

// Notifier 是客户端期望的接口
type Notifier interface {
	Send(message string)
}

// EmailNotifier 是具体的通知组件
type EmailNotifier struct{}

func (e *EmailNotifier) Send(message string) {
	fmt.Printf("Sending email notification: %s\n", message)
}

// SMSNotifier 是具体的通知组件
type SMSNotifier struct{}

func (s *SMSNotifier) Send(message string) {
	fmt.Printf("Sending SMS notification: %s\n", message)
}

// Decorator 是装饰器接口
type Decorator interface {
	Notifier
}

// UrgentNotifier 是装饰器，用于添加紧急通知功能
type UrgentNotifier struct {
	notifier Notifier
}

func (u *UrgentNotifier) Send(message string) {
	// Add urgent functionality here
	fmt.Printf("Sending URGENT notification: %s\n", message)
	u.notifier.Send(message)
}

func main() {
	// 创建原始通知组件
	emailNotifier := &EmailNotifier{}
	smsNotifier := &SMSNotifier{}

	// 添加紧急通知功能
	urgentEmailNotifier := &UrgentNotifier{notifier: emailNotifier}
	urgentSMSNotifier := &UrgentNotifier{notifier: smsNotifier}

	// 发送通知
	urgentEmailNotifier.Send("Payment received")
	urgentSMSNotifier.Send("Package delivered")
}
