package fileflow

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Dispatcher struct {
	processors     []Processor
	middlewares    []Middleware
	writer         *ResultWriter
	errHandler     ErrorHandler
	sem            *semaphore.Weighted
	retryPolicy    RetryPolicy
	processTimeout time.Duration
}

// Dispatch 消费事件并按并发上限调度处理，退出前会等待在途任务完成。
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

// handle 执行单文件处理流程：匹配 Processor -> 中间件链 -> 重试 -> 写结果/错误处理。
func (d *Dispatcher) handle(ctx context.Context, event FileEvent) {
	proc := d.match(event)
	if proc == nil {
		return
	}

	chain := d.buildChain(proc)

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
		d.errHandler.Handle(ctx, event, err)
		return
	}

	result.SourceEvent = event
	if werr := d.writer.Write(ctx, result); werr != nil {
		d.errHandler.Handle(ctx, event, werr)
	}
}

// match 按注册顺序找到首个可处理该文件的 Processor。
func (d *Dispatcher) match(e FileEvent) Processor {
	for _, p := range d.processors {
		if p.Match(e) {
			return p
		}
	}
	return nil
}

// waitContext 在可取消上下文中等待 backoff 时间。
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
