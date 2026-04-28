package gitutil

import (
	"sync"

	"github.com/nef734/ggit/internal/conf"
)

var (
	poolOnce sync.Once
	pool     chan struct{}
)

func getPool() chan struct{} {
	poolOnce.Do(func() {
		limit := conf.Git.MaxConcurrency
		if limit <= 0 {
			limit = 50
		}
		pool = make(chan struct{}, limit)
		for i := 0; i < limit; i++ {
			pool <- struct{}{}
		}
	})
	return pool
}

// Acquire gets a token from the pool.
func Acquire() {
	<-getPool()
}

// Release returns a token to the pool.
func Release() {
	getPool() <- struct{}{}
}
