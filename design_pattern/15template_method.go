package main

import "fmt"

// AbstractClass defines the template method
type AbstractClass interface {
	TemplateMethod()
}

// ConcreteClassA implements the AbstractClass interface
type ConcreteClassA struct{}

func (c *ConcreteClassA) TemplateMethod() {
	fmt.Println("Executing strategy A")
}

// ConcreteClassB implements the AbstractClass interface
type ConcreteClassB struct{}

func (c *ConcreteClassB) TemplateMethod() {
	fmt.Println("Executing strategy B")
}

func main() {
	var strategyA, strategyB AbstractClass
	// Create concrete strategy instances
	strategyA = &ConcreteClassA{}
	strategyB = &ConcreteClassB{}

	// Use strategy A
	strategyA.TemplateMethod()

	// Switch to strategy B
	strategyA = strategyB
	strategyA.TemplateMethod()
}
