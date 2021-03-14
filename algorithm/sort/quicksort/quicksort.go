package main

import "fmt"

func main() {
	c := []int{2, 6, 4, 8, 5, 3}
	QuickSort1(c, 0, len(c)-1)
	println(fmt.Sprintf("%v", c))
}

func QuickSort1(arr []int, b, e int) {
	if b >= e {
		return
	}
	m := Partition1(arr, b, e)
	QuickSort1(arr, b, m-1)
	QuickSort1(arr, m+1, e)
}

func Partition1(arr []int, b, e int) int {
	sentry := arr[b]
	l, r := b, e
	rTurn := true
	for l != r {
		if rTurn {
			if arr[r] < sentry {
				arr[l], arr[r] = arr[r], arr[l]
				rTurn = false
			}
		} else {
			if arr[l] > sentry {
				arr[l], arr[r] = arr[r], arr[l]
				rTurn = true
			}
		}
		if rTurn {
			r--
		} else {
			l++
		}
	}
	return l
}

func QuickSort2(arr []int, p, r int) {
	if p < r {
		q := Partition2(arr, p, r)
		QuickSort2(arr, p, q-1)
		QuickSort2(arr, q+1, r)
	}
}

func Partition2(arr []int, p, r int) int {
	x := arr[r]
	i := p - 1
	for j := p; j <= r-1; j++ {
		if arr[j] <= x {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[r] = arr[r], arr[i+1]
	return i + 1
}
