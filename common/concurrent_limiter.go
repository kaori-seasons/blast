package common

import (
	"errors"
	"time"
)

type ConcurrentLimiter struct {
	limiter chan bool
}

func NewConcurrentLimiter(max uint32) *ConcurrentLimiter {
	r := &ConcurrentLimiter{limiter: make(chan bool, max)}
	for i := uint32(0); i < max; i++ {
		r.limiter <- true
	}

	return r
}

func (c *ConcurrentLimiter) Acquire(timeout_ms time.Duration) (func(), error) {
	select {
	case <-c.limiter:
		return func() { c.limiter <- true }, nil
	case <-time.After(timeout_ms * time.Millisecond):
		return nil, errors.New("time out")
	}
}
