//go:build linux

package osutil

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// AvailableMemory returns the available memory in bytes.
func AvailableMemory() (uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return 0, err
			}
			return val * 1024, nil // /proc/meminfo is in kB
		}
	}
	return 0, scanner.Err()
}
