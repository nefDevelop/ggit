//go:build windows

package osutil

import (
	"syscall"
	"unsafe"
)

var (
	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

// AvailableMemory returns the available physical memory in bytes.
func AvailableMemory() (uint64, error) {
	var msx memoryStatusEx
	msx.cbSize = uint32(unsafe.Sizeof(msx))
	r1, _, e1 := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&msx)))
	if r1 == 0 {
		if e1 != nil {
			return 0, e1
		}
		return 0, syscall.EINVAL
	}
	return msx.ullAvailPhys, nil
}
