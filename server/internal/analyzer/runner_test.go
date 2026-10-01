package analyzer

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/urlvet/urlvet/internal/service/threatfeeds"
)

type fastTask struct{}

func (fastTask) Name() string { return "fast" }
func (fastTask) Run(_ *Input, out *Output) error {
	updateOutputDirect(out, func(o *Output) { o.Rank = 7 })
	return nil
}

// stuckTask stands in for a check against a site that never answers.
type stuckTask struct{ release chan struct{} }

func (stuckTask) Name() string { return "stuck" }
func (s stuckTask) Run(_ *Input, out *Output) error {
	<-s.release
	updateOutputDirect(out, func(o *Output) { o.Rank = 99 })
	return nil
}

func TestRunTasks_ReturnsAtDeadline(t *testing.T) {
	release := make(chan struct{})
	defer close(release)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	out, errs := runTasks(ctx, &Input{}, []Task{fastTask{}, stuckTask{release}})
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("runTasks waited %v for a stuck task", elapsed)
	}
	if out.Rank != 7 {
		t.Errorf("Rank = %d, want the finished task's 7", out.Rank)
	}
	if len(errs) == 0 || !strings.Contains(errs[len(errs)-1].Error(), "stuck") {
		t.Errorf("errs = %v, want a timeout naming the stuck task", errs)
	}
}

func TestSummarizeErrors(t *testing.T) {
	rateLimited := &TaskError{Task: "phishtank_check", Err: threatfeeds.ErrRateLimited}
	unreachable := &TaskError{Task: "content_check", Err: errors.New("dial tcp: i/o timeout")}
	timeout := &TimeoutError{Tasks: []string{"content_check", "tls_combined_check"}}

	incomplete, checks, msgs := summarizeErrors([]error{rateLimited})
	if incomplete {
		t.Error("a PhishTank rate limit alone should not make the result incomplete")
	}
	if !slices.Equal(checks, []string{"phishtank_check"}) || len(msgs) != 1 {
		t.Errorf("checks = %v, msgs = %v", checks, msgs)
	}

	incomplete, checks, _ = summarizeErrors([]error{rateLimited, unreachable, timeout})
	if !incomplete {
		t.Error("an unreachable site should make the result incomplete")
	}
	if want := []string{"content_check", "phishtank_check", "tls_combined_check"}; !slices.Equal(checks, want) {
		t.Errorf("checks = %v, want %v", checks, want)
	}
}
