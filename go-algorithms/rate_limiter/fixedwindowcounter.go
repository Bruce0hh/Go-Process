package rate_limiter

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	slots map[int]int // 时间窗口中的计数器
	mu    sync.Mutex  // 互斥锁，保证并发安全
}

func NewFixedWindowCounter() *FixedWindowCounter {
	return &FixedWindowCounter{
		slots: make(map[int]int),
	}
}

// Incr 请求进来，增加计数器
func (fwc *FixedWindowCounter) Incr() {
	now := time.Now().Unix()
	fwc.mu.Lock()
	fwc.slots[int(now)]++
	fwc.mu.Unlock()
}

// Count 统计最近一分钟内计数器的总数，并清除过期的计数器
func (fwc *FixedWindowCounter) Count() int {
	now := time.Now().Unix()
	count := 0
	fwc.mu.Lock()
	for timestamp, c := range fwc.slots {
		if int(now)-timestamp < 60 { // 只统计最近一分钟的计数器
			count += c
		} else {
			delete(fwc.slots, timestamp) // 清除过期的计数器
		}
	}
	fwc.mu.Unlock()
	return count
}
