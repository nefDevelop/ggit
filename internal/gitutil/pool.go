package gitutil

import (
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/nef734/ggit/internal/conf"
)

var (
	poolOnce sync.Once
	pool     chan struct{}
)

func calculateRAMLimit() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		// Fallback to CPU cores if RAM check fails
		return runtime.NumCPU() * 2
	}

	// Blindaje Anti-OOM Heuristic: 
	// Each git process can spike up to 256MB-512MB in heavy operations (like diff/pack).
	// We use 50% of available RAM as a safe concurrency zone.
	availableMB := v.Available / 1024 / 1024
	limit := int(availableMB / 384) // Conservative 384MB per process

	if limit < 2 {
		limit = 2 // Minimum for basic operations
	}
	if limit > 64 {
		limit = 64 // Hard cap to prevent thread explosion
	}
	return limit
}

func getPool() chan struct{} {
	poolOnce.Do(func() {
		limit := conf.Git.MaxConcurrency
		if limit <= 0 {
			limit = calculateRAMLimit()
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
