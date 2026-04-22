package tool

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// List of global semaphores to limit concurrent heavy processes.
var (
	// GitGrepSemaphore limits the number of concurrent 'git grep' processes.
	// This prevents memory exhaustion under high load.
	GitGrepSemaphore = semaphore.NewWeighted(10) // Limit to 10 concurrent greps
)

// AcquireSemaphore is a helper to acquire a semaphore with context.
func AcquireSemaphore(ctx context.Context, s *semaphore.Weighted) error {
	return s.Acquire(ctx, 1)
}

// ReleaseSemaphore is a helper to release a semaphore.
func ReleaseSemaphore(s *semaphore.Weighted) {
	s.Release(1)
}
