/**
滑动窗口限流算法，基于ring环形链表
*/

package ring

import (
	"container/ring"
	"sync"
	"sync/atomic"
	"time"
)

var defaultSlidingWindow = SlidingWindow{
	size:     3,
	bucket:   3,
	interval: 1000,
}

// SlidingWindow 滑动窗口
type SlidingWindow struct {
	size     int64 // 窗口期内最大请求数
	bucket   int   // 滑动窗口个数
	count    int64 // 所有窗口的累计请求数
	interval int64 // 窗口周期ms
	head     *ring.Ring
	mu       sync.Mutex
}

type Option func(*SlidingWindow)

// NewSlidingWindow returns SlidingWindow object
func NewSlidingWindow(opts ...Option) *SlidingWindow {
	p := &defaultSlidingWindow
	for _, opt := range opts {
		opt(p)
	}

	head := ring.New(p.bucket)
	for i := 0; i < p.bucket; i++ {
		head.Value = 0
		head = head.Next()
	}

	sw := &SlidingWindow{
		size:     p.size,
		bucket:   p.bucket,
		interval: p.interval,
		head:     head,
	}

	go sw.ticker()

	return sw
}

func (sw *SlidingWindow) ticker() {
	// 启动执行器，每隔interval毫秒刷新一次滑动窗口数据
	ticker := time.NewTicker(time.Duration(sw.interval) * time.Millisecond)
	for range ticker.C {
		atomic.AddInt64(&sw.count, int64(0-sw.head.Value.(int)))
		sw.head.Value = 0
		sw.head = sw.head.Next()
	}
}

// Exec 执行
func (sw *SlidingWindow) Exec() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.count >= sw.size {
		return false
	}

	atomic.AddInt64(&sw.count, 1)
	pos := sw.head.Prev()
	pos.Value = pos.Value.(int) + 1

	return true
}

// WithSize 设置size
func WithSize(size int64) Option {
	return func(window *SlidingWindow) {
		window.size = size
	}
}

// WithBucket 设置窗口数
func WithBucket(bucket int) Option {
	return func(window *SlidingWindow) {
		window.bucket = bucket
	}
}

// WithInterval 设置interval
func WithInterval(interval int64) Option {
	return func(window *SlidingWindow) {
		window.interval = interval
	}
}
