package redis

import (
	"log"
	"testing"
	"time"
)

func TestPeriodLimit(t *testing.T) {
	p := NewPeriodLimit("test-period-key")
	for i := 0; i < 5; i++ {
		n := 0
		switch i {
		case 0:
			n = 2
		case 1:
			n = 20
		case 2:
			n = 8
		case 3:
			n = 5
		case 4:
			n = 6
		}
		for j := 0; j < n; j++ {
			resp, err := p.Exec()
			if err != nil {
				log.Printf("exec error: %v\n", err)
				continue
			}

			switch resp {
			case LimitAllowed:
				log.Println("allow")
			case LimitHitQuota:
				log.Println("hit")
			case LimitOverQuota:
				log.Println("over")
			}
		}

		time.Sleep(time.Second)
	}

}
