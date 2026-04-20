package initialize

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/fileflow"
	"go.uber.org/zap"
)

var fileflowRuntime struct {
	mu     sync.Mutex
	cancel context.CancelFunc
}

// Fileflow 按配置启动 fileflow 管道（默认关闭）。
func Fileflow() {
	cfg := global.GVA_CONFIG.Fileflow
	if !cfg.Enable {
		return
	}

	// 拿到业务侧已注册的处理器与扩展点。
	processors, middlewares, postProcessor := fileflow.SnapshotRegistered()
	if len(processors) == 0 {
		zap.L().Warn("fileflow is enabled but no processor was registered")
		return
	}

	ffCfg, err := buildFileflowConfig(cfg)
	if err != nil {
		zap.L().Error("failed to build fileflow config", zap.Error(err))
		return
	}

	pipeline, err := fileflow.New(ffCfg)
	if err != nil {
		zap.L().Error("failed to initialize fileflow pipeline", zap.Error(err))
		return
	}

	// 先挂内置中间件，再挂业务自定义中间件。
	pipeline.
		UseMiddleware(fileflow.RecoverMiddleware()).
		UseMiddleware(fileflow.LoggingMiddleware(ffCfg.Logger))

	for _, m := range middlewares {
		pipeline.UseMiddleware(m)
	}
	for _, p := range processors {
		pipeline.Use(p)
	}
	if postProcessor != nil {
		pipeline.WithPostProcessor(postProcessor)
	}

	startFileflow(pipeline)
}

// StopFileflow 停止当前运行中的 fileflow 管道。
func StopFileflow() {
	fileflowRuntime.mu.Lock()
	defer fileflowRuntime.mu.Unlock()
	if fileflowRuntime.cancel != nil {
		fileflowRuntime.cancel()
		fileflowRuntime.cancel = nil
	}
}

// startFileflow 在独立协程中启动管道。
func startFileflow(p *fileflow.Pipeline) {
	StopFileflow()

	fileflowRuntime.mu.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	fileflowRuntime.cancel = cancel
	fileflowRuntime.mu.Unlock()

	go func() {
		err := p.Run(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			zap.L().Error("fileflow exited with error", zap.Error(err))
		}
	}()
}

// buildFileflowConfig 将 YAML 配置转换为运行时 Config。
func buildFileflowConfig(cfg config.Fileflow) (fileflow.Config, error) {
	if strings.TrimSpace(cfg.WatchDir) == "" {
		return fileflow.Config{}, errors.New("fileflow: watch-dir is required when enabled")
	}
	if strings.TrimSpace(cfg.OutputDir) == "" {
		return fileflow.Config{}, errors.New("fileflow: output-dir is required when enabled")
	}

	interval, err := parseDurationOrDefault(cfg.Interval, time.Second)
	if err != nil {
		return fileflow.Config{}, fmt.Errorf("parse fileflow interval: %w", err)
	}
	timeout, err := parseDurationOrDefault(cfg.Timeout, 30*time.Second)
	if err != nil {
		return fileflow.Config{}, fmt.Errorf("parse fileflow timeout: %w", err)
	}
	backoffBase, err := parseDurationOrDefault(cfg.BackoffBase, 500*time.Millisecond)
	if err != nil {
		return fileflow.Config{}, fmt.Errorf("parse fileflow backoff base: %w", err)
	}
	afterProcess, err := parseAfterProcess(cfg.AfterProcess)
	if err != nil {
		return fileflow.Config{}, err
	}

	// 默认不过滤；ignore-hidden=true 时过滤隐藏文件。
	filter := func(name string) bool { return true }
	if cfg.IgnoreHidden {
		filter = func(name string) bool { return !strings.HasPrefix(name, ".") }
	}

	ffCfg := fileflow.Config{
		WatchDir:       filepath.Clean(cfg.WatchDir),
		OutputDir:      filepath.Clean(cfg.OutputDir),
		Interval:       interval,
		Concurrency:    cfg.Concurrency,
		EventBuffer:    cfg.EventBuffer,
		ProcessTimeout: timeout,
		Filter:         filter,
		AfterProcess:   afterProcess,
		RetryPolicy: fileflow.RetryPolicy{
			MaxRetries: cfg.MaxRetries,
			Backoff:    fileflow.ExponentialBackoff(backoffBase),
		},
		Logger: slog.Default(),
	}
	if strings.TrimSpace(cfg.FailedDir) != "" {
		ffCfg.FailedDir = filepath.Clean(cfg.FailedDir)
	}
	return ffCfg, nil
}

// parseDurationOrDefault 解析 duration 字符串，空值时回退默认值。
func parseDurationOrDefault(raw string, defaultValue time.Duration) (time.Duration, error) {
	if strings.TrimSpace(raw) == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(raw)
}

// parseAfterProcess 解析源文件后处理策略。
func parseAfterProcess(raw string) (fileflow.AfterWriteAction, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "move", "move-to-processed":
		return fileflow.MoveToProcessed, nil
	case "delete":
		return fileflow.DeleteSource, nil
	case "keep":
		return fileflow.KeepSource, nil
	default:
		return 0, fmt.Errorf("fileflow: unsupported after-process action %q", raw)
	}
}
