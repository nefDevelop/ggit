package osx

import (
	"os"

	"github.com/nef734/ggit/internal/errx"
)

var _ errx.NotFound = (*Error)(nil)

// Error is a wrapper of an OS error, which handles not found.
type Error struct {
	error
}

func (e Error) NotFound() bool {
	return e.error == os.ErrNotExist
}

// NewError wraps given error.
func NewError(err error) error {
	return Error{error: err}
}
