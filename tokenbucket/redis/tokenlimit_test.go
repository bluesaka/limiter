package redis

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestTokenLimit(t *testing.T) {
	tl := NewTokenLimiter("test-token-limit", WithRate(5), WithCapacity(10))
	wg := sync.WaitGroup{}
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if tl.AllowN(time.Now(), 1) {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
	}

	wg.Wait()
}
