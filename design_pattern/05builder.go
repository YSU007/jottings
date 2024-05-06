package main

import "fmt"

// PersonInterface 接口表示一个人的信息和行为
type PersonInterface interface {
	Speak(message string)
	Sleep()
}

// Person 结构体实现了 PersonInterface 接口
type Person struct {
	name    string
	age     int
	gender  string
	address string
}

func (p *Person) Speak(message string) {
	fmt.Printf("%s says: %s\n", p.name, message)
}

func (p *Person) Sleep() {
	fmt.Printf("%s is sleeping now...\n", p.name)
}

// PersonBuilderInterface 接口表示一个人的信息构造器
type PersonBuilderInterface interface {
	SetName(name string) PersonBuilderInterface
	SetAge(age int) PersonBuilderInterface
	SetGender(gender string) PersonBuilderInterface
	SetAddress(address string) PersonBuilderInterface
	Build() PersonInterface
}

// PersonBuilder 结构体表示一个人的建造器
type PersonBuilder struct {
	name    string
	age     int
	gender  string
	address string
}

func (b *PersonBuilder) SetName(name string) PersonBuilderInterface {
	b.name = name
	return b
}

func (b *PersonBuilder) SetAge(age int) PersonBuilderInterface {
	b.age = age
	return b
}

func (b *PersonBuilder) SetGender(gender string) PersonBuilderInterface {
	b.gender = gender
	return b
}

func (b *PersonBuilder) SetAddress(address string) PersonBuilderInterface {
	b.address = address
	return b
}

func (b *PersonBuilder) Build() PersonInterface {
	return &Person{
		name:    b.name,
		age:     b.age,
		gender:  b.gender,
		address: b.address,
	}
}

func main() {
	// 创建一个 PersonBuilder 实例
	builder := &PersonBuilder{}

	// 设置人的信息
	person := builder.SetName("Tom").
		SetAge(18).
		SetGender("Male").
		SetAddress("123 Main St").
		Build()

	// 调用人的说话和睡觉方法
	person.Speak("Hello, World!")
	person.Sleep()
}
