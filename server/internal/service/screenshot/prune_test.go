package screenshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPruneScreenshots(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "screenshot-old.png")
	fresh := filepath.Join(dir, "screenshot-fresh.png")
	for _, f := range []string{old, fresh} {
		if err := os.WriteFile(f, []byte("png"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	past := time.Now().Add(-25 * time.Hour)
	if err := os.Chtimes(old, past, past); err != nil {
		t.Fatal(err)
	}

	if n := PruneScreenshots(dir, 24*time.Hour); n != 1 {
		t.Errorf("removed %d files, want 1", n)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("expired screenshot was not deleted")
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Error("fresh screenshot was deleted")
	}
	if n := PruneScreenshots(filepath.Join(dir, "missing"), time.Hour); n != 0 {
		t.Errorf("missing dir: removed %d, want 0", n)
	}
}
