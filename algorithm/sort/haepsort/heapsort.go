package main

import (
	"fmt"
)

func main() {
	//no use index 0
	c := []int{0, 2, 6, 4, 8, 5, 3}
	HeapSort(c)
	println(fmt.Sprintf("%v", c))
}

func Parent(i int) int {
	return i / 2
}

func Left(i int) int {
	return 2 * i
}

func Right(i int) int {
	return 2*i + 1
}

func MaxHeapify(arr []int, i, heapSize int) {
	largest := i
	l := Left(i)
	r := Right(i)
	if l <= heapSize && arr[l] > arr[largest] {
		largest = l
	}
	if r <= heapSize && arr[r] > arr[largest] {
		largest = r
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		MaxHeapify(arr, largest, heapSize)
	}
}

func BuildMaxHeap(arr []int, heapSize int) {
	for i := heapSize / 2; i >= 1; i-- {
		MaxHeapify(arr, i, heapSize)
	}
}

func HeapSort(arr []int) {
	heapSize := len(arr) - 1
	for i := len(arr) - 1; i >= 2; i-- {
		BuildMaxHeap(arr, heapSize)
		arr[1], arr[i] = arr[i], arr[1]
		heapSize--
	}
}
