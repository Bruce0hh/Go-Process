package rate_limiter

import (
	"sync"
	"time"
)

type LeakingBucket struct {
	capacity   int        // 桶容量
	waterLevel int        // 当前水位
	outflow    int        // 每秒出流量
	mu         sync.Mutex // 互斥锁，保证并发安全
}

func NewLeakingBucket(capacity int, outflow int) *LeakingBucket {
	lb := &LeakingBucket{
		capacity:   capacity,
		waterLevel: capacity,
		outflow:    outflow,
	}
	go lb.startLeak()
	return lb
}

// 开始漏一个outflow的水，实际上是控制消费的速率
func (lb *LeakingBucket) startLeak() {
	for {
		time.Sleep(time.Second)
		lb.mu.Lock()
		lb.waterLevel = lb.waterLevel - lb.outflow
		if lb.waterLevel < 0 {
			lb.waterLevel = 0
		}
		lb.mu.Unlock()
	}
}

// PourWater 加水
func (lb *LeakingBucket) PourWater(amount int) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	// 加水后水位下降
	if lb.waterLevel >= amount {
		lb.waterLevel = lb.waterLevel - amount
		return true
	}
	return false
}
