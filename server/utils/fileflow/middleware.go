package fileflow

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"
)

type ProcessFunc func(ctx context.Context, e FileEvent) (Result, error)
type Middleware func(next ProcessFunc) ProcessFunc

// buildChain 按“先注册先执行”的顺序构建中间件链。
func (d *Dispatcher) buildChain(proc Processor) ProcessFunc {
	h := ProcessFunc(proc.Process)
	for i := len(d.middlewares) - 1; i >= 0; i-- {
		h = d.middlewares[i](h)
	}
	return h
}

// RecoverMiddleware 兜底捕获 panic，转成 error 返回。
func RecoverMiddleware() Middleware {
	return func(next ProcessFunc) ProcessFunc {
		return func(ctx context.Context, e FileEvent) (r Result, err error) {
			defer func() {
				if rec := recover(); rec != nil {
					err = fmt.Errorf("panic: %v\n%s", rec, debug.Stack())
				}
			}()
			return next(ctx, e)
		}
	}
}

// LoggingMiddleware 记录单文件处理耗时与结果。
func LoggingMiddleware(logger *slog.Logger) Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	return func(next ProcessFunc) ProcessFunc {
		return func(ctx context.Context, e FileEvent) (Result, error) {
			start := time.Now()
			r, err := next(ctx, e)
			logger.Info(
				"fileflow: processed file",
				"file", e.Name,
				"duration_ms", time.Since(start).Milliseconds(),
				"error", err,
			)
			return r, err
		}
	}
}

// TimeoutMiddleware 为单文件处理创建独立超时上下文。
func TimeoutMiddleware(d time.Duration) Middleware {
	return func(next ProcessFunc) ProcessFunc {
		return func(ctx context.Context, e FileEvent) (Result, error) {
			tctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(tctx, e)
		}
	}
}
