package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

var wg = sync.WaitGroup{}

func busi(ch chan int) {
	for t := range ch {
		fmt.Println("go task = ", t, ", goroutine count = ", runtime.NumGoroutine())
		wg.Done()
	}
}

func sendTask(task int, ch chan int) {
	wg.Add(1)
	ch <- task
}

func main() {
	ch := make(chan int) // 无buffer channel

	goCnt := 3 // 启动goroutine的数量
	for i := 0; i < goCnt; i++ {
		// 启动go
		go busi(ch)
	}

	taskCnt := math.MaxInt64 // 模拟用户需求业务的数量
	for t := 0; t < taskCnt; t++ {
		// 发送任务
		sendTask(t, ch)
	}

	wg.Wait()
}
