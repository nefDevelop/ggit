package mocks

import (
	"gogs.io/gogs/internal/context"
)

var _ context.Locale = (*Locale)(nil)

type Locale struct {
	MockTr func(string, ...any) string
}

func (l *Locale) Tr(format string, args ...any) string {
	return l.MockTr(format, args...)
}
