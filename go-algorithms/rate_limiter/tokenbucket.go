package rate_limiter

import (
	"sync"
	"time"
)

// TokenBucket 可根据速率和容量实现不同的桶，配置不同的限流规则
type TokenBucket struct {
	capacity   int           // 桶容量
	tokens     int           // 当前令牌数量
	refillRate time.Duration // 令牌补充速率
	mu         sync.Mutex    // 互斥锁，保证并发安全
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
	}
	go tb.startRefill()
	return tb
}

// 补充token
func (tb *TokenBucket) startRefill() {
	for {
		time.Sleep(tb.refillRate)
		tb.mu.Lock()
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
		tb.mu.Unlock()
	}
}

// TakeToken 获取token
func (tb *TokenBucket) TakeToken() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}
