package fileflow

import (
	"errors"
	"log/slog"
	"path/filepath"
	"time"
)

type Config struct {
	// WatchDir 为输入目录，必须配置。
	WatchDir string
	// OutputDir 为输出目录，必须配置。
	OutputDir string

	// FailedDir 为失败文件目录，默认 WatchDir/failed。
	FailedDir string
	// Interval 为扫描间隔。
	Interval time.Duration
	// Concurrency 为并发处理上限。
	Concurrency int
	// EventBuffer 为 watcher 事件通道缓冲区大小。
	EventBuffer int
	// ProcessTimeout 为单文件处理超时时间。
	ProcessTimeout time.Duration
	// Filter 用于按文件名过滤输入文件。
	Filter func(string) bool
	// AfterProcess 控制处理成功后的源文件动作。
	AfterProcess AfterWriteAction
	// RetryPolicy 控制处理失败后的重试策略。
	RetryPolicy RetryPolicy
	// Logger 为日志实现，未配置时使用 slog.Default。
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
)

type RetryPolicy struct {
	MaxRetries int
	Backoff    func(attempt int) time.Duration
}

// ExponentialBackoff 返回指数退避函数。
func ExponentialBackoff(base time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		if attempt < 1 {
			return base
		}
		return base * time.Duration(1<<uint(attempt-1))
	}
}

// defaultConfig 补齐默认值并校验关键字段。
func defaultConfig(cfg Config) (Config, error) {
	if cfg.WatchDir == "" {
		return cfg, errors.New("fileflow: watch dir is required")
	}
	if cfg.OutputDir == "" {
		return cfg, errors.New("fileflow: output dir is required")
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
