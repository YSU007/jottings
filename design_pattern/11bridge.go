package main

import "fmt"

// Shape 表示图形接口
type Shape interface {
	Draw()
}

// Color 表示颜色接口
type Color interface {
	Fill()
}

// Circle 是具体的图形
type Circle struct {
	color Color
}

func (c *Circle) Draw() {
	fmt.Println("Drawing a circle")
	c.color.Fill()
}

// RedColor 是具体的颜色
type RedColor struct{}

func (r *RedColor) Fill() {
	fmt.Println("Filling with red color")
}

// BlueColor 是具体的颜色
type BlueColor struct{}

func (b *BlueColor) Fill() {
	fmt.Println("Filling with blue color")
}

func main() {
	redCircle := &Circle{color: &RedColor{}}
	blueCircle := &Circle{color: &BlueColor{}}

	redCircle.Draw()
	blueCircle.Draw()
}
