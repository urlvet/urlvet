package logger

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestRedactURLs(t *testing.T) {
	tests := map[string]string{
		"https://a.com/reset?token=abc123":                    "https://a.com/reset?[redacted]",
		"http://x.io/p#frag":                                  "http://x.io/p?[redacted]",
		"https://example.com/plain":                           "https://example.com/plain",
		`Get "https://bank.com/login?sid=9&u=me": timeout`:    `Get "https://bank.com/login?[redacted]": timeout`,
		"no url here":                                         "no url here",
		"two https://a.com/?x=1 and https://b.com/y?z=2 done": "two https://a.com/?[redacted] and https://b.com/y?[redacted] done",
	}
	for in, want := range tests {
		if got := RedactURLs(in); got != want {
			t.Errorf("RedactURLs(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRedactingHandler(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(&redactingHandler{next: slog.NewTextHandler(&buf, nil)})
	log.With("base", "https://w.com/a?k=v").Info(
		"fetching https://m.com/?q=secret",
		"url", "https://s.com/x?token=t1",
		"err", errors.New(`Get "https://e.com/r?code=c2": EOF`),
		slog.Group("g", "inner", "https://g.com/?p=s3"),
	)
	out := buf.String()
	for _, secret := range []string{"k=v", "q=secret", "token=t1", "code=c2", "p=s3"} {
		if strings.Contains(out, secret) {
			t.Errorf("log output leaked %q: %s", secret, out)
		}
	}
	if !strings.Contains(out, "https://s.com/x?[redacted]") {
		t.Errorf("expected redacted URL to keep host and path: %s", out)
	}
}
