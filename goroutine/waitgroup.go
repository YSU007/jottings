package main

import (
	"sync"
)

const N = 10

var wg1 = &sync.WaitGroup{}

func main() {
	for i := 0; i < N; i++ {
		wg1.Add(1)
		go func(i int) {
			println(i)
			defer wg1.Done()
		}(i)
	}

	wg.Wait()
}
