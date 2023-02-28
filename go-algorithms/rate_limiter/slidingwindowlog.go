package rate_limiter

import (
	"sync"
	"time"
)

type SlidingWindowLog struct {
	logs []int64    // 时间窗口中的请求时间戳
	size int        // 时间窗口的大小
	pos  int        // 滑动窗口的起始位置
	mu   sync.Mutex // 互斥锁，保证并发安全
}

func NewSlidingWindowLog(size int) *SlidingWindowLog {
	return &SlidingWindowLog{
		logs: make([]int64, size),
		size: size,
		pos:  0,
	}
}

// AddLog 添加日志，更新pos位置
func (swl *SlidingWindowLog) AddLog(timestamp int64) {
	swl.mu.Lock()
	swl.logs[swl.pos] = timestamp
	swl.pos = (swl.pos + 1) % swl.size // 更新滑动窗口的位置
	swl.mu.Unlock()
}

// Count 统计最近一分钟的请求总数
func (swl *SlidingWindowLog) Count() int {
	now := time.Now().Unix()
	count := 0
	swl.mu.Lock()
	for i := 0; i < swl.size; i++ {
		if now-swl.logs[i] < 60 { // 只统计最近一分钟的请求
			count++
		}
	}
	swl.mu.Unlock()
	return count
}
