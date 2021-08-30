package ring

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestSlidingWindow(t *testing.T) {
	sw := NewSlidingWindow(WithSize(3), WithInterval(1000))
	wg := sync.WaitGroup{}
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if sw.Exec() {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
		time.Sleep(time.Millisecond * 600)
	}

	wg.Wait()
}
