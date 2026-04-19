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

// New 创建并初始化 fileflow 流水线。
func New(cfg Config) (*Pipeline, error) {
	var err error
	cfg, err = defaultConfig(cfg)
	if err != nil {
		return nil, err
	}

	writer := &ResultWriter{
		outputDir:  cfg.OutputDir,
		afterWrite: cfg.AfterProcess,
		doneDir:    filepath.Join(cfg.WatchDir, "done"),
	}

	dispatcher := &Dispatcher{
		writer: writer,
		sem:    semaphore.NewWeighted(int64(cfg.Concurrency)),
		errHandler: &DefaultErrorHandler{
			FailedDir: cfg.FailedDir,
			Logger:    cfg.Logger,
		},
		retryPolicy:    cfg.RetryPolicy,
		processTimeout: cfg.ProcessTimeout,
	}

	return &Pipeline{
		cfg:        cfg,
		watcher:    newFileWatcher(cfg),
		dispatcher: dispatcher,
	}, nil
}

// Use 注册业务 Processor（按注册顺序匹配）。
func (p *Pipeline) Use(proc Processor) *Pipeline {
	p.dispatcher.processors = append(p.dispatcher.processors, proc)
	return p
}

// UseMiddleware 注册中间件。
func (p *Pipeline) UseMiddleware(m Middleware) *Pipeline {
	p.dispatcher.middlewares = append(p.dispatcher.middlewares, m)
	return p
}

// WithPostProcessor 设置可选后置处理器。
func (p *Pipeline) WithPostProcessor(pp PostProcessor) *Pipeline {
	p.dispatcher.writer.postProcessor = pp
	return p
}

// WithErrorHandler 自定义错误处理器。
func (p *Pipeline) WithErrorHandler(errHandler ErrorHandler) *Pipeline {
	if errHandler != nil {
		p.dispatcher.errHandler = errHandler
	}
	return p
}

// Run 启动流水线并阻塞到 ctx 取消或事件消费结束。
func (p *Pipeline) Run(ctx context.Context) error {
	if len(p.dispatcher.processors) == 0 {
		return errors.New("fileflow: no processors registered")
	}
	events := p.watcher.Watch(ctx)
	p.dispatcher.Dispatch(ctx, events)
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil
	}
	return ctx.Err()
}
