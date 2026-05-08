//go:build darwin

package osutil

import (
	"syscall"
	"unsafe"
)

// AvailableMemory returns the available memory in bytes for macOS.
func AvailableMemory() (uint64, error) {
	// For macOS, we'll use sysctl to get page size and free page count
	pageSize, err := syscall.SysctlUint32("vm.pagesize")
	if err != nil {
		return 0, err
	}

	freeCount, err := syscall.SysctlUint32("vm.page_free_count")
	if err != nil {
		return 0, err
	}

	// Also consider inactive pages as "available"
	inactiveCount, err := syscall.SysctlUint32("vm.page_inactive_count")
	if err != nil {
		// If we can't get inactive, just use free
		return uint64(freeCount) * uint64(pageSize), nil
	}

	return uint64(freeCount+inactiveCount) * uint64(pageSize), nil
}
