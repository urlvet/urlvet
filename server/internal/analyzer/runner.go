package analyzer

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/urlvet/urlvet/internal/metrics"
	"github.com/urlvet/urlvet/internal/store"
)

// TaskError is a task that failed.
type TaskError struct {
	Task string
	Err  error
}

func (e *TaskError) Error() string { return e.Task + ": " + e.Err.Error() }
func (e *TaskError) Unwrap() error { return e.Err }

// TimeoutError lists the tasks still running when the scan's deadline passed.
type TimeoutError struct{ Tasks []string }

func (e *TimeoutError) Error() string { return fmt.Sprintf("timed out waiting for: %v", e.Tasks) }

// runTasks runs all tasks in parallel, collects timings and non-fatal errors.
//
// Tasks don't take a context, so one stuck on a slow site (some block automated
// visitors by never answering) would hold up the whole scan. When ctx expires,
// runTasks stops waiting and returns a copy of what has finished, with an error
// naming the tasks that hadn't. Late tasks keep writing to the original Output,
// which nothing reads any more.
func runTasks(ctx context.Context, in *Input, tasks []Task) (*Output, []error) {
	out := &Output{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	finished := map[string]bool{}

	for _, t := range tasks {
		t := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					msg := fmt.Sprintf("%s: panic: %v", t.Name(), r)
					mu.Lock()
					errs = append(errs, &TaskError{Task: t.Name(), Err: fmt.Errorf("panic: %v", r)})
					mu.Unlock()
					metrics.TaskErrors.WithLabelValues(t.Name()).Inc()
					store.AddError(store.ErrorRecord{
						Task:  t.Name(),
						Error: msg,
						URL:   in.URL,
						Time:  time.Now(),
					})
				}
			}()
			select {
			case <-ctx.Done():
				return
			default:
				start := time.Now()
				err := t.Run(in, out)
				elapsed := time.Since(start)
				metrics.TaskDuration.WithLabelValues(t.Name()).Observe(elapsed.Seconds())
				if err != nil {
					metrics.TaskErrors.WithLabelValues(t.Name()).Inc()
					mu.Lock()
					errs = append(errs, &TaskError{Task: t.Name(), Err: err})
					mu.Unlock()
				}
				out.setTiming(t.Name(), elapsed)
				mu.Lock()
				finished[t.Name()] = true
				mu.Unlock()
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return out, errs
	case <-ctx.Done():
	}

	out.mu.Lock()
	snapshot := &Output{OutputData: out.OutputData}
	snapshot.Timings = maps.Clone(out.Timings)
	out.mu.Unlock()

	mu.Lock()
	defer mu.Unlock()
	var pending []string
	for _, t := range tasks {
		if !finished[t.Name()] {
			pending = append(pending, t.Name())
		}
	}
	sort.Strings(pending)
	return snapshot, append(slices.Clone(errs), &TimeoutError{Tasks: pending})
}
