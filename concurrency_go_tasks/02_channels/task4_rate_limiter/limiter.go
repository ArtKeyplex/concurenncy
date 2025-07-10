package limiter

import (
	"time"
)

type Limiter struct {
	tokens chan struct{}
}

func NewLimiter() *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, 5),
	}

	for range 5 {
		l.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for i := len(l.tokens); i < 5; i++ {
				select {
				case l.tokens <- struct{}{}:
				default:
					break
				}
			}
		}
	}()

	return l
}

func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens: 
		return true
	default: 
		return false
	}
}
