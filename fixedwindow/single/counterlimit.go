/**
固定窗口算法
*/

package single

import (
	"sync"
	"time"
)

var defaultCounter = Counter{
	quota:  10,
	period: time.Second,
}

// Counter 计数器
type Counter struct {
	quota  int           // 计数周期内最大请求数
	count  int           // 计数周期内累计请求数
	period time.Duration // 计数周期
	begin  time.Time     // 计数开始时间
	mu     sync.Mutex
}

type Option func(*Counter)

// NewCounter returns Counter object
func NewCounter(opts ...Option) *Counter {
	p := &defaultCounter
	for _, opt := range opts {
		opt(p)
	}
	return &Counter{
		quota:  p.quota,
		period: p.period,
		begin:  time.Now(),
	}
}

// Exec 执行计数器
func (c *Counter) Exec() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.count == c.quota {
		if time.Since(c.begin) >= c.period {
			c.Reset()
			return true
		} else {
			return false
		}
	} else {
		c.count++
		return true
	}
}

// Reset 重置计数器
func (c *Counter) Reset() {
	c.begin = time.Now()
	c.count = 0
}

// WithQuota 设置计数周期内的最大请求数
func WithQuota(quota int) Option {
	return func(c *Counter) {
		c.quota = quota
	}
}

// WithPeriod 设置计数周期
func WithPeriod(period time.Duration) Option {
	return func(c *Counter) {
		c.period = period
	}
}
