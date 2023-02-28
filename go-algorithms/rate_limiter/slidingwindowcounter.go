package rate_limiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	limits int         // 单位窗口的限流数
	size   int         // 单位窗口的时间窗格数
	pos    int         // 单位窗口的末位
	Slots  map[int]int // 时间窗格中的请求数
	mu     sync.Mutex  // 互斥锁，保证并发安全
}

func NewSlidingWindowCounter(limits, size int) *SlidingWindowCounter {
	swc := &SlidingWindowCounter{
		limits: limits,
		Slots:  make(map[int]int, size),
		size:   size,
		pos:    int(time.Now().Unix()),
	}
	go swc.sliding()
	return swc
}

// 移动一个时间窗口的距离
func (swc *SlidingWindowCounter) sliding() {
	for {
		time.Sleep(time.Second) // 假设一个时间窗格是1s
		swc.mu.Lock()
		now := int(time.Now().Unix())
		swc.pos = now
		delete(swc.Slots, now-swc.size*1)
		swc.mu.Unlock()
	}
}

func (swc *SlidingWindowCounter) Incr() {
	if swc.limits <= swc.Count() {
		return
	}
	swc.mu.Lock()
	swc.Slots[swc.pos]++
	swc.mu.Unlock()
}

func (swc *SlidingWindowCounter) Count() int {
	count := 0
	for _, v := range swc.Slots {
		count += v
	}
	return count
}
