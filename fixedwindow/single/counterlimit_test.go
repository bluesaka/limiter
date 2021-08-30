package single

import (
	"log"
	"sync"
	"testing"
)

func TestCounterLimit(t *testing.T) {
	// 每秒最多5个请求
	counter := NewCounter(WithQuota(5))
	wg := sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if counter.Exec() {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
	}

	wg.Wait()
}
