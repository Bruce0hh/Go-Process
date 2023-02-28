package main

import (
	"fmt"
	"go_algorithms/rate_limiter"
	"time"
)

func main() {
	// 创建一个每10秒限流10次的滑动窗口计数器
	swc := rate_limiter.NewSlidingWindowCounter(10, 10)

	// 每200毫秒调用一次 Incr
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			swc.Incr()
			fmt.Printf("Current count: %d\n", swc.Count())
		}
	}

	//滑动日志
	//swl := rate_limiter.NewSlidingWindowLog(60)
	//for i := 0; i < 100; i++ {
	//	swl.AddLog(time.Now().Unix())
	//	fmt.Printf("Requests in last minute: %d\n", swl.Count())
	//	time.Sleep(time.Second)
	//}
}
