package single

import (
	"log"
	"sync"
	"testing"
)

func TestLeakyBucket(t *testing.T) {
	lb := NewLeakyBucket(WithRate(5), WithCapacity(10))
	wg := sync.WaitGroup{}
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if lb.Exec() {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
	}

	wg.Wait()
}

func TestLeakyBucket2(t *testing.T) {
	lb := NewLeakyBucket(WithRate(5), WithCapacity(10))
	for i := 1; i <= 20; i++ {
		if lb.Exec() {
			log.Printf("i: %d allow", i)
		} else {
			log.Printf("i: %d over", i)
		}
	}

}
