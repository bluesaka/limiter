/**
golang标准库自带的令牌桶限流算法
golang.org/x/time/rate
 */

package main

import (
	"context"
	"golang.org/x/time/rate"
	"log"
)

func main() {
	// 每秒向桶中放入4个token，桶容量为2
	limiter := rate.NewLimiter(4, 2)
	for i := 1; i <= 10; i++ {
		// Wait/WaitN
		// Allow/AllowN
		// Reserve/ReserveN
		if err := limiter.WaitN(context.Background(), 1); err != nil {
			log.Println(err)
		} else {
			log.Printf("i: %d pass", i)
		}
	}
}
