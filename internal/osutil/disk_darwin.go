//go:build darwin

package osutil

import (
	"syscall"
)

// DiskUsage returns total and free bytes for the given path.
func DiskUsage(path string) (total, free uint64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}
	return stat.Blocks * uint64(stat.Bsize), stat.Bfree * uint64(stat.Bsize), nil
}
