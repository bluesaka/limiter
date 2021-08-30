/**
滑动窗口限流算法
*/

package single

import (
	"sync"
	"time"
)

var defaultSlidingWindow = SlidingWindow{
	size:     2,
	interval: 1000,
}

// SlidingWindow 滑动窗口
type SlidingWindow struct {
	key      string
	size     int64   // 窗口期内最大请求数
	interval int64   // 窗口周期ms
	window   []int64 // 窗口请求
	mu       sync.Mutex
}

type Option func(*SlidingWindow)

// NewSlidingWindow returns SlidingWindow object
func NewSlidingWindow(key string, opts ...Option) *SlidingWindow {
	p := &defaultSlidingWindow
	for _, opt := range opts {
		opt(p)
	}
	return &SlidingWindow{
		key:      key,
		size:     p.size,
		interval: p.interval,
		window:   make([]int64, 0, p.size),
	}
}

// Exec 执行
func (sw *SlidingWindow) Exec() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // 当前毫秒数
	if int64(len(sw.window)) < sw.size {
		sw.window = append(sw.window, now)
		return true
	}

	head := sw.window[0]
	// 最早时间的请求还在时间窗口期内，拒绝此次请求
	if now-head <= sw.interval {
		return false
	}

	// 最早时间的请求在时间窗口期外，去除并接受此次请求
	sw.window = append(sw.window[1:], now)
	return true
}

// WithSize 设置size
func WithSize(size int64) Option {
	return func(window *SlidingWindow) {
		window.size = size
	}
}

// WithInterval 设置interval
func WithInterval(interval int64) Option {
	return func(window *SlidingWindow) {
		window.interval = interval
	}
}
