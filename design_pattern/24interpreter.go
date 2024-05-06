package main

import (
	"fmt"
)

// Expression 是表达式接口
type Expression interface {
	Interpret(variables map[string]Expression) int
}

// Number 是一个表达式实现，代表数字
type Number struct {
	value int
}

func (n *Number) Interpret(variables map[string]Expression) int {
	return n.value
}

// Plus 是一个表达式实现，代表加法
type Plus struct {
	leftOperand  Expression
	rightOperand Expression
}

func (p *Plus) Interpret(variables map[string]Expression) int {
	return p.leftOperand.Interpret(variables) + p.rightOperand.Interpret(variables)
}

// Minus 是一个表达式实现，代表减法
type Minus struct {
	leftOperand  Expression
	rightOperand Expression
}

func (m *Minus) Interpret(variables map[string]Expression) int {
	return m.leftOperand.Interpret(variables) - m.rightOperand.Interpret(variables)
}

// Variable 是一个表达式实现，代表变量
type Variable struct {
	name string
}

func (v *Variable) Interpret(variables map[string]Expression) int {
	if val, ok := variables[v.name]; ok {
		return val.Interpret(variables)
	}
	return 0
}

func main() {
	expression := &Plus{
		leftOperand: &Minus{
			leftOperand:  &Number{value: 7},
			rightOperand: &Number{value: 3},
		},
		rightOperand: &Plus{
			leftOperand:  &Variable{name: "x"},
			rightOperand: &Number{value: 1},
		},
	}

	variables := make(map[string]Expression)
	variables["x"] = &Number{value: 2}

	fmt.Printf("The result of the expression is: %d\n", expression.Interpret(variables))
}
