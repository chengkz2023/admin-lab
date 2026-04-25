package fileflow

import (
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	// WatchDir 为轮询输入目录。
	WatchDir string
	// FailedDir 存放重试后仍失败的源文件。
	FailedDir string
	// Interval 为目录轮询间隔。
	Interval time.Duration
	// Concurrency 限制并发处理数量。
	Concurrency int
	// EventBuffer 为 watcher 事件通道缓冲大小。
	EventBuffer int
	// ProcessTimeout 限制单次处理调用耗时。
	ProcessTimeout time.Duration
	// Filter 判断文件名是否进入处理候选集。
	Filter func(string) bool
	// ReadyStrategy 判断发现的文件是否已就绪可处理。
	ReadyStrategy ReadyStrategy
	// AfterProcess 控制默认 Hook 在处理成功后的源文件动作。
	AfterProcess AfterWriteAction
	// RetryPolicy 控制处理失败后的重试策略。
	RetryPolicy RetryPolicy
	// Logger 接收框架日志。
	Logger *slog.Logger
}

type AfterWriteAction int

const (
	// MoveToProcessed 将源文件移动到 WatchDir/done。
	MoveToProcessed AfterWriteAction = iota
	// DeleteSource 删除源文件。
	DeleteSource
	// KeepSource 保留源文件。
	KeepSource
	// ArchiveOnly 是 MoveToProcessed 的语义别名，适用于只归档源文件的处理器。
	ArchiveOnly
)

type RetryPolicy struct {
	MaxRetries int
	Backoff    func(attempt int) time.Duration
}

// ExponentialBackoff 返回指数退避重试间隔函数。
func ExponentialBackoff(base time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		if attempt < 1 {
			return base
		}
		return base * time.Duration(1<<uint(attempt-1))
	}
}

// ExtFilter 返回按一个或多个扩展名过滤文件名的函数。
func ExtFilter(exts ...string) func(string) bool {
	allowed := make(map[string]struct{}, len(exts))
	for _, ext := range exts {
		if strings.TrimSpace(ext) == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		allowed[strings.ToLower(ext)] = struct{}{}
	}
	return func(name string) bool {
		_, ok := allowed[strings.ToLower(filepath.Ext(name))]
		return ok
	}
}

// AndFilter 用逻辑与组合多个过滤器，nil 过滤器会被忽略。
func AndFilter(filters ...func(string) bool) func(string) bool {
	return func(name string) bool {
		for _, filter := range filters {
			if filter != nil && !filter(name) {
				return false
			}
		}
		return true
	}
}

// defaultConfig 补齐默认值并校验必要字段。
func defaultConfig(cfg Config) (Config, error) {
	if cfg.WatchDir == "" {
		return cfg, errors.New("fileflow: watch dir is required")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 4
	}
	if cfg.EventBuffer <= 0 {
		cfg.EventBuffer = 64
	}
	if cfg.ProcessTimeout <= 0 {
		cfg.ProcessTimeout = 30 * time.Second
	}
	if cfg.FailedDir == "" {
		cfg.FailedDir = filepath.Join(cfg.WatchDir, "failed")
	}
	if cfg.ReadyStrategy == nil {
		cfg.ReadyStrategy = RenameReady{}
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.RetryPolicy.MaxRetries < 0 {
		return cfg, errors.New("fileflow: retry max retries cannot be negative")
	}
	if cfg.RetryPolicy.Backoff == nil {
		cfg.RetryPolicy.Backoff = ExponentialBackoff(500 * time.Millisecond)
	}
	return cfg, nil
}
