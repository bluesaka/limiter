/**
令牌桶算法
*/

package single

import (
	"sync"
	"time"
)

var defaultLeakyBucket = LeakyBucket{
	rate:     5,
	capacity: 10,
}

// LeakyBucket 漏桶
type LeakyBucket struct {
	rate      int   //固定每秒出水速率
	capacity  int   // 桶容量
	water     int   // 桶中当前水量
	timestamp int64 // 桶上次出水毫秒时间戳
	mu        sync.Mutex
}

type Option func(*LeakyBucket)

// NewLeakyBucket returns LeakyBucket object
func NewLeakyBucket(opts ...Option) *LeakyBucket {
	p := &defaultLeakyBucket
	for _, opt := range opts {
		opt(p)
	}
	return &LeakyBucket{
		rate:     p.rate,
		capacity: p.capacity,
	}
}

// Exec 执行
func (lb *LeakyBucket) Exec() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	lb.water = max(0, lb.water-int(now-lb.timestamp)*lb.rate/1000)
	lb.timestamp = now

	if lb.water < lb.capacity {
		lb.water++
		return true
	} else {
		return false
	}
}

// WithRate 设置漏桶速率
func WithRate(rate int) Option {
	return func(lb *LeakyBucket) {
		lb.rate = rate
	}
}

// WithCapacity 设置漏桶容量
func WithCapacity(capacity int) Option {
	return func(lb *LeakyBucket) {
		lb.capacity = capacity
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
