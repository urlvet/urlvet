package middleware

import (
	"strings"
	"testing"
	"time"
)

func TestIPHasherKeys(t *testing.T) {
	h := newIPHasher()

	a := h.key("203.0.113.7")
	if strings.Contains(a, "203.0.113.7") {
		t.Fatalf("key leaks the IP: %s", a)
	}
	if !strings.HasPrefix(a, "ratelimit:") {
		t.Fatalf("missing prefix: %s", a)
	}
	if h.key("203.0.113.7") != a {
		t.Fatal("same IP should map to the same key within a rotation")
	}
	if h.key("203.0.113.8") == a {
		t.Fatal("different IPs should map to different keys")
	}
}

func TestIPHasherRotates(t *testing.T) {
	now := time.Unix(0, 0)
	h := &ipHasher{now: func() time.Time { return now }}

	before := h.key("203.0.113.7")
	now = now.Add(saltRotation)
	if h.key("203.0.113.7") == before {
		t.Fatal("key should change once the salt rotates")
	}
}
