package logger

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
)

// Scanned links can carry secrets in their query string (password-reset tokens,
// invite codes, session IDs). Logs must never store them, so every record passes
// through redactingHandler, which strips query strings and fragments from URLs
// found in the message, string attributes and error values.

var urlPattern = regexp.MustCompile(`https?://[^\s"'<>]+`)

// RedactURLs removes the query string and fragment from every URL in s,
// keeping scheme, host and path: "https://a.com/reset?token=x" → "https://a.com/reset?[redacted]".
func RedactURLs(s string) string {
	return urlPattern.ReplaceAllStringFunc(s, func(u string) string {
		if i := strings.IndexAny(u, "?#"); i >= 0 {
			return u[:i] + "?[redacted]"
		}
		return u
	})
}

type redactingHandler struct {
	next slog.Handler
}

func (h *redactingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *redactingHandler) Handle(ctx context.Context, r slog.Record) error {
	out := slog.NewRecord(r.Time, r.Level, RedactURLs(r.Message), r.PC)
	r.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(redactAttr(a))
		return true
	})
	return h.next.Handle(ctx, out)
}

func (h *redactingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clean := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		clean[i] = redactAttr(a)
	}
	return &redactingHandler{next: h.next.WithAttrs(clean)}
}

func (h *redactingHandler) WithGroup(name string) slog.Handler {
	return &redactingHandler{next: h.next.WithGroup(name)}
}

func redactAttr(a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindString:
		return slog.String(a.Key, RedactURLs(a.Value.String()))
	case slog.KindAny:
		if err, ok := a.Value.Any().(error); ok && err != nil {
			return slog.String(a.Key, RedactURLs(err.Error()))
		}
	case slog.KindGroup:
		group := a.Value.Group()
		clean := make([]any, len(group))
		for i, g := range group {
			clean[i] = redactAttr(g)
		}
		return slog.Group(a.Key, clean...)
	}
	return a
}
