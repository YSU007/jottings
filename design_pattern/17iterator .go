package main

import "fmt"

// Iterator 接口定义迭代器方法
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

// Ints 是一个整数切片
type Ints []int

// NewInts 创建一个整数切片
func NewInts(data ...int) Ints {
	return data
}

// Iterator 返回整数切片的迭代器
func (i Ints) Iterator() Iterator {
	return &IntsIterator{data: i, index: 0}
}

// IntsIterator 是整数切片的迭代器
type IntsIterator struct {
	data  Ints
	index int
}

// HasNext 判断是否还有下一个元素
func (it *IntsIterator) HasNext() bool {
	return it.index < len(it.data)
}

// Next 返回下一个元素
func (it *IntsIterator) Next() interface{} {
	if it.HasNext() {
		val := it.data[it.index]
		it.index++
		return val
	}
	return nil
}

func main() {
	ints := NewInts(1, 2, 3, 4, 5)
	iterator := ints.Iterator()

	for iterator.HasNext() {
		fmt.Println(iterator.Next())
	}
}
