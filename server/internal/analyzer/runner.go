package analyzer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/urlvet/urlvet/internal/metrics"
	"github.com/urlvet/urlvet/internal/store"
)

// runTasks runs all tasks in parallel, collects timings and non-fatal errors
func runTasks(ctx context.Context, in *Input, tasks []Task) (*Output, []error) {
	out := &Output{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, t := range tasks {
		t := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					msg := fmt.Sprintf("%s: panic: %v", t.Name(), r)
					mu.Lock()
					errs = append(errs, fmt.Errorf("%s", msg))
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
					errs = append(errs, fmt.Errorf("%s: %w", t.Name(), err))
					mu.Unlock()
				}
				out.setTiming(t.Name(), elapsed)
			}
		}()
	}

	wg.Wait()
	return out, errs
}
