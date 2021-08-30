/**
漏桶算法

uber的ratelimiter算法是漏桶算法的一种实现
 */

package main

import (
	"go.uber.org/ratelimit"
	"log"
	"time"
)

func main() {
	rl := ratelimit.New(2) // 每秒处理2个请求
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		log.Println(i, now.Sub(prev))
		prev = now
	}
}
