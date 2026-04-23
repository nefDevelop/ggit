package mocks

import (
	"github.com/nef734/ggit/internal/context"
)

var _ context.Locale = (*Locale)(nil)

type Locale struct {
	MockTr       func(string, ...any) string
	MockLanguage func() string
}

func (l *Locale) Tr(format string, args ...any) string {
	return l.MockTr(format, args...)
}

func (l *Locale) Language() string {
	if l.MockLanguage != nil {
		return l.MockLanguage()
	}
	return "en-US"
}
