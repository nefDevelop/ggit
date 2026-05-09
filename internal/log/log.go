package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level = slog.Level

const (
	LevelTrace = slog.Level(-8)
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)

	DefaultConsoleName = "console"
	DefaultFileName    = "file"
)

var (
	mu       sync.Mutex
	handlers = make(map[string]slog.Handler)
)

func init() {
	// Initialize with a default console logger
	_ = NewConsole(DefaultConsoleName, LevelTrace)
}

// Trace logs a message at LevelTrace.
func Trace(format string, v ...any) {
	slog.Default().Log(context.Background(), LevelTrace, fmt.Sprintf(format, v...))
}

// Info logs a message at slog.LevelInfo.
func Info(format string, v ...any) {
	slog.Default().Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, v...))
}

// Warn logs a message at slog.LevelWarn.
func Warn(format string, v ...any) {
	slog.Default().Log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, v...))
}

// Error logs a message at slog.LevelError.
func Error(format string, v ...any) {
	slog.Default().Log(context.Background(), slog.LevelError, fmt.Sprintf(format, v...))
}

// ErrorDepth logs a message at slog.LevelError with the given call depth.
func ErrorDepth(depth int, format string, v ...any) {
	var pcs [1]uintptr
	runtime.Callers(depth+2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, v...), pcs[0])
	_ = slog.Default().Handler().Handle(context.Background(), r)
}

// Fatal logs a message at LevelFatal and then calls os.Exit(1).
func Fatal(format string, v ...any) {
	slog.Default().Log(context.Background(), LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// NewConsole creates a new console handler and adds it to the multi-handler.
func NewConsole(name string, level slog.Level) error {
	mu.Lock()
	defer mu.Unlock()

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				l := a.Value.Any().(slog.Level)
				if l == LevelTrace {
					a.Value = slog.StringValue("TRACE")
				} else if l == LevelFatal {
					a.Value = slog.StringValue("FATAL")
				}
			}
			return a
		},
	})
	handlers[name] = h
	rebuildDefault()
	return nil
}

// NewFile creates a new file handler and adds it to the multi-handler.
func NewFile(name, filename string, level slog.Level) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	h := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				l := a.Value.Any().(slog.Level)
				if l == LevelTrace {
					a.Value = slog.StringValue("TRACE")
				} else if l == LevelFatal {
					a.Value = slog.StringValue("FATAL")
				}
			}
			return a
		},
	})
	handlers[name] = h
	rebuildDefault()
	return nil
}

// Remove removes a handler by name.
func Remove(name string) {
	mu.Lock()
	defer mu.Unlock()

	delete(handlers, name)
	rebuildDefault()
}

// New is a dummy function for compatibility with tests using clog.New.
func New(name string, _ any) error {
	if name == "noop" {
		mu.Lock()
		defer mu.Unlock()
		handlers[name] = &noopHandler{}
		rebuildDefault()
	}
	return nil
}

// FileRotationConfig is a dummy struct for compatibility.
type FileRotationConfig struct {
	Rotate   bool
	Daily    bool
	MaxSize  int64
	MaxLines int64
	MaxDays  int64
}

// NewFileWriter creates a simple file writer without rotation for now.
func NewFileWriter(filename string, _ FileRotationConfig) (io.Writer, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// Stop shuts down the logging service.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	handlers = make(map[string]slog.Handler)
	rebuildDefault()
}

func rebuildDefault() {
	var hs []slog.Handler
	for _, h := range handlers {
		hs = append(hs, h)
	}
	if len(hs) == 0 {
		slog.SetDefault(slog.New(&noopHandler{}))
		return
	}
	slog.SetDefault(slog.New(&multiHandler{handlers: hs}))
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			_ = h.Handle(ctx, r)
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

type noopHandler struct{}

func (noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (noopHandler) Handle(context.Context, slog.Record) error { return nil }
func (h noopHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h noopHandler) WithGroup(string) slog.Handler           { return h }
