package fileflow

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Dispatcher struct {
	processor      Processor
	middlewares    []Middleware
	hook           Hook
	sem            *semaphore.Weighted
	retryPolicy    RetryPolicy
	processTimeout time.Duration
}

// Dispatch 消费事件，并按并发上限调度处理。
func (d *Dispatcher) Dispatch(ctx context.Context, events <-chan FileEvent) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			if err := d.sem.Acquire(ctx, 1); err != nil {
				return
			}
			wg.Add(1)
			go func(e FileEvent) {
				defer wg.Done()
				defer d.sem.Release(1)
				d.handle(ctx, e)
			}(event)
		}
	}
}

func (d *Dispatcher) handle(ctx context.Context, event FileEvent) {
	if d.processor == nil {
		return
	}

	chain := d.buildChain(d.processor)
	var result Result
	var err error

	for attempt := 0; attempt <= d.retryPolicy.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := d.retryPolicy.Backoff(attempt)
			if !waitContext(ctx, backoff) {
				return
			}
		}

		procCtx := ctx
		cancel := func() {}
		if d.processTimeout > 0 {
			procCtx, cancel = context.WithTimeout(ctx, d.processTimeout)
		}
		result, err = chain(procCtx, event)
		cancel()

		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return
		}
	}

	if err != nil {
		if d.hook != nil {
			d.hook.OnError(ctx, event, err)
		}
		return
	}

	if d.hook != nil {
		if hookErr := d.hook.OnSuccess(ctx, event, result); hookErr != nil {
			d.hook.OnError(ctx, event, hookErr)
		}
	}
}

func waitContext(ctx context.Context, delay time.Duration) bool {
	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}
