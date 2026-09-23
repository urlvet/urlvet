// Package logger provides a centralized structured logger built on log/slog.
//
// Call Init() once after environment variables are loaded. In DEV mode it
// emits colored human-readable text to stderr and plain text to a rotating
// daily log file; in production it emits JSON to both. Timezone is controlled
// by LOG_TIMEZONE (e.g. "Asia/Kolkata"); defaults to UTC. Log files are
// written under LOG_DIR (e.g. "logs"); leave empty to disable file logging.
package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// ANSI color codes (DEV terminal only)
// ---------------------------------------------------------------------------

const (
	colorReset  = "\033[0m"
	colorDim    = "\033[2m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[97m"
)

// ---------------------------------------------------------------------------
// rotatingFileWriter — io.Writer that rotates to a new file each calendar day
// ---------------------------------------------------------------------------

type rotatingFileWriter struct {
	mu      sync.Mutex
	dir     string
	loc     *time.Location
	current string
	file    *os.File
}

func (w *rotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().In(w.loc)
	path := filepath.Join(w.dir, now.Format("2006-01"), now.Format("02")+".log")

	if path != w.current {
		if w.file != nil {
			w.file.Close()
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return 0, err
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.current = path
	}

	return w.file.Write(p)
}

func (w *rotatingFileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		w.file.Close()
		w.file = nil
	}
}

// ---------------------------------------------------------------------------
// multiHandler — fans a slog.Record out to multiple handlers
// ---------------------------------------------------------------------------

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var firstErr error
	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, r.Level) {
			continue
		}
		if err := handler.Handle(ctx, r.Clone()); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		out[i] = hh.WithAttrs(attrs)
	}
	return &multiHandler{handlers: out}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	out := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		out[i] = hh.WithGroup(name)
	}
	return &multiHandler{handlers: out}
}

// ---------------------------------------------------------------------------
// colorHandler — colored human-readable output for DEV terminal
// ---------------------------------------------------------------------------

type colorHandler struct {
	out   io.Writer
	level slog.Level
	loc   *time.Location
	attrs []slog.Attr
}

func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	// Timestamp — dimmed so it doesn't compete with the level badge
	ts := r.Time.In(h.loc).Format("2006-01-02 15:04:05 MST")
	fmt.Fprintf(&buf, "%s%s%s ", colorDim, ts, colorReset)

	// Level badge, padded to 7 chars for alignment
	lc := levelColor(r.Level)
	fmt.Fprintf(&buf, "%s%-7s%s ", lc+colorBold, "["+r.Level.String()+"]", colorReset)

	// Message
	fmt.Fprintf(&buf, "%s%s%s", colorWhite+colorBold, r.Message, colorReset)

	for _, a := range h.attrs {
		writeColorAttr(&buf, a)
	}
	r.Attrs(func(a slog.Attr) bool {
		writeColorAttr(&buf, a)
		return true
	})

	buf.WriteByte('\n')
	_, err := h.out.Write(buf.Bytes())
	return err
}

func writeColorAttr(buf *bytes.Buffer, a slog.Attr) {
	fmt.Fprintf(buf, " %s%s%s=%s%v%s", colorDim, a.Key, colorReset, colorCyan, a.Value.Any(), colorReset)
}

func levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan
	}
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(merged, h.attrs)
	copy(merged[len(h.attrs):], attrs)
	return &colorHandler{out: h.out, level: h.level, loc: h.loc, attrs: merged}
}

func (h *colorHandler) WithGroup(name string) slog.Handler { return h }

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

// Init configures the global logger based on ENV, LOG_TIMEZONE, and LOG_DIR.
// Must be called after godotenv.Load() so env vars are already set.
func Init() {
	loc := loadTimezone()
	logDir := os.Getenv("LOG_DIR")

	replaceTime := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().In(loc).Format(time.RFC3339))
		}
		return a
	}

	var handler slog.Handler

	if os.Getenv("ENV") == "DEV" {
		termHandler := &colorHandler{
			out:   os.Stderr,
			level: slog.LevelDebug,
			loc:   loc,
		}

		if logDir == "" {
			handler = termHandler
		} else {
			fw := &rotatingFileWriter{dir: logDir, loc: loc}
			fileHandler := slog.NewTextHandler(fw, &slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Value = slog.StringValue(a.Value.Time().In(loc).Format("2006-01-02 15:04:05 MST"))
					}
					return a
				},
			})
			handler = &multiHandler{handlers: []slog.Handler{termHandler, fileHandler}}
		}
	} else {
		var w io.Writer = os.Stderr
		if logDir != "" {
			w = io.MultiWriter(os.Stderr, &rotatingFileWriter{dir: logDir, loc: loc})
		}
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: replaceTime,
		})
	}

	// Strip query strings from any URL before it reaches a log (see redact.go).
	slog.SetDefault(slog.New(&redactingHandler{next: handler}))
}

// loadTimezone reads LOG_TIMEZONE and returns the matching *time.Location.
// Falls back to UTC on unknown values.
func loadTimezone() *time.Location {
	tz := os.Getenv("LOG_TIMEZONE")
	if tz == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: unknown timezone %q, falling back to UTC\n", tz)
		return time.UTC
	}
	return loc
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }

// Fatal logs at Error level then exits. Mirrors log.Fatal behaviour.
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
