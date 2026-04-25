package fileflow

import (
	"context"
	"errors"
	"path/filepath"

	"golang.org/x/sync/semaphore"
)

type Pipeline struct {
	cfg        Config
	watcher    *FileWatcher
	dispatcher *Dispatcher
}

// New 创建 fileflow 流水线。
func New(cfg Config) (*Pipeline, error) {
	var err error
	cfg, err = defaultConfig(cfg)
	if err != nil {
		return nil, err
	}

	defaultHook := &DefaultHook{
		DoneDir:    filepath.Join(cfg.WatchDir, "done"),
		FailedDir:  cfg.FailedDir,
		AfterWrite: cfg.AfterProcess,
		Logger:     cfg.Logger,
	}
	dispatcher := &Dispatcher{
		hook:           defaultHook,
		sem:            semaphore.NewWeighted(int64(cfg.Concurrency)),
		retryPolicy:    cfg.RetryPolicy,
		processTimeout: cfg.ProcessTimeout,
	}

	return &Pipeline{
		cfg:        cfg,
		watcher:    newFileWatcher(cfg),
		dispatcher: dispatcher,
	}, nil
}

// Use 注册单个业务 Processor，后续调用会覆盖前一次注册。
func (p *Pipeline) Use(proc Processor) *Pipeline {
	p.dispatcher.processor = proc
	return p
}

// UseMiddleware 注册包裹 Processor.Process 的中间件。
func (p *Pipeline) UseMiddleware(m Middleware) *Pipeline {
	if m != nil {
		p.dispatcher.middlewares = append(p.dispatcher.middlewares, m)
	}
	return p
}

// WithHook 替换流水线 Hook。
func (p *Pipeline) WithHook(hook Hook) *Pipeline {
	if hook != nil {
		p.dispatcher.hook = hook
	}
	return p
}

// Run 启动流水线，并阻塞到 ctx 取消或事件消费完成。
func (p *Pipeline) Run(ctx context.Context) error {
	if p.dispatcher.processor == nil {
		return errors.New("fileflow: no processor registered")
	}
	events := p.watcher.Watch(ctx)
	p.dispatcher.Dispatch(ctx, events)
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil
	}
	return ctx.Err()
}
