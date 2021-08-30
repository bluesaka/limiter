package single

import (
	"log"
	"sync"
	"testing"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(WithRate(5), WithCapacity(10))
	wg := sync.WaitGroup{}
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if tb.Exec() {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
	}

	wg.Wait()
}
