package threatfeeds

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"os"
	"testing"
)

// The browser extension normalizes links the same way before hashing them;
// web/chrome-extension tests read this same file.
func TestFeedKeyCases(t *testing.T) {
	raw, err := os.ReadFile("testdata/feedkey_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ URL, Key string }
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		if got := FeedKey(c.URL); got != c.Key {
			t.Errorf("FeedKey(%q) = %q, want %q", c.URL, got, c.Key)
		}
	}
}

func TestExportPrefixes(t *testing.T) {
	saved := local
	defer func() { local = saved }()
	local = &store{urls: map[string]map[string]struct{}{}, hosts: map[string]map[string]struct{}{}}
	exportGen = -1

	local.set("PhishTank", []string{"https://evil.example/login", "http://www.evil.example/login/"}) // same key twice
	local.set("OpenPhish", []string{"https://not-redistributable.example/"})

	exp := ExportPrefixes()
	if exp.Count != 1 || len(exp.Prefixes) != PrefixLen {
		t.Fatalf("count %d, %d bytes; want 1 prefix (deduplicated, OpenPhish left out)", exp.Count, len(exp.Prefixes))
	}
	sum := sha256.Sum256([]byte("evil.example/login"))
	if !bytes.Equal(exp.Prefixes, sum[:PrefixLen]) {
		t.Errorf("prefix %x, want %x", exp.Prefixes, sum[:PrefixLen])
	}
	if len(exp.Sources) != 1 || exp.Sources[0] != "PhishTank" {
		t.Errorf("sources %v, want only PhishTank", exp.Sources)
	}

	// Cached until a feed changes; a change gives a new version.
	if again := ExportPrefixes(); again.Version != exp.Version {
		t.Error("version changed without a feed change")
	}
	local.set("URLhaus", []string{"http://1.2.3.4/x"})
	if next := ExportPrefixes(); next.Count != 2 || next.Version == exp.Version {
		t.Errorf("after URLhaus loaded: count %d, version changed %v", next.Count, next.Version != exp.Version)
	}
}
