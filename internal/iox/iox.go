package iox

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies the file at src to dst, preserving file mode and
// modification time.
func CopyFile(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source"+": %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source"+": %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create target"+": %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy"+": %w", err)
	}

	if err = out.Sync(); err != nil {
		return fmt.Errorf("sync target"+": %w", err)
	}

	if err = os.Chmod(dst, si.Mode()); err != nil {
		return fmt.Errorf("chmod target"+": %w", err)
	}

	if err = os.Chtimes(dst, si.ModTime(), si.ModTime()); err != nil {
		return fmt.Errorf("chtimes target"+": %w", err)
	}

	return nil
}
