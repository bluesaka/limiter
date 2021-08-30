package single

import (
	"sync"
	"time"
)

var defaultTokenBucket = TokenBucket{
	rate:     2000,
	capacity: 2000,
}

// TokenBucket 令牌桶
type TokenBucket struct {
	rate      int64 // 每秒放入token数
	capacity  int64 // 令牌桶容量
	token     int64 // 当前桶内token数
	timestamp int64 // 最后一次取token的时间戳
	mu        sync.Mutex
}

type Option func(*TokenBucket)

// NewTokenBucket returns TokenBucket object
func NewTokenBucket(opts ...Option) *TokenBucket {
	p := &defaultTokenBucket
	for _, opt := range opts {
		opt(p)
	}
	return &TokenBucket{
		rate:      p.rate,
		capacity:  p.capacity,
		timestamp: time.Now().Unix() - 1,
	}
}

// Exec 执行
func (tb *TokenBucket) Exec() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now().Unix()
	// 添加令牌
	tb.token = tb.token + (now-tb.timestamp)*tb.rate
	if tb.token > tb.capacity {
		tb.token = tb.capacity
	}
	tb.timestamp = now

	if tb.token > 0 {
		tb.token--
		return true
	} else {
		return false
	}
}

// WithRate 设置令牌桶每秒放入token数
func WithRate(rate int64) Option {
	return func(lb *TokenBucket) {
		lb.rate = rate
	}
}

// WithCapacity 设置令牌桶容量
func WithCapacity(capacity int64) Option {
	return func(lb *TokenBucket) {
		lb.capacity = capacity
	}
}
