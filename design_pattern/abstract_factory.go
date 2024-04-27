package main

import "fmt"

// Food 是抽象的食物接口
type Food interface {
	Eat()
}

// Drug 是抽象的药物接口
type Drug interface {
	Take()
}

// Factory 负责生产药物类和食物类
type Factory interface {
	NewFood() Food
	NewDrug() Drug
}

// 实现食物类结构体
type meat struct {
	// 肉
}

func (m meat) Eat() {
	fmt.Println("Eat meat")
}

type fruit struct {
	// 水果
}

func (f fruit) Eat() {
	fmt.Println("Eat fruit")
}

// 实现药物类结构体
type feverdrug struct {
	// 发烧药
}

func (f feverdrug) Take() {
	fmt.Println("Take fever drug")
}

type colddrug struct {
	// 感冒药
}

func (c colddrug) Take() {
	fmt.Println("Take cold drug")
}

// 第一个工厂负责肉和发烧药的生成
type FirstFactory struct {
	// 第一个工厂
}

func (f *FirstFactory) NewFood() Food {
	return meat{}
}

func (f *FirstFactory) NewDrug() Drug {
	return feverdrug{}
}

// 第二个工厂负责水果和感冒药的生成
type SecondFactory struct {
	// 第二个工厂
}

func (f *SecondFactory) NewFood() Food {
	return fruit{}
}

func (f *SecondFactory) NewDrug() Drug {
	return colddrug{}
}

func main() {
	fmt.Println("抽象工厂模式")

	first := &FirstFactory{}
	first.NewFood().Eat()
	first.NewDrug().Take()

	second := &SecondFactory{}
	second.NewFood().Eat()
	second.NewDrug().Take()
}
