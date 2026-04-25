package initialize

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/fileflow"
	"go.uber.org/zap"
)

// Fileflow 是项目内 fileflow 流水线的统一装配入口。
//
// 设计约定：一个目录对应一个独立 Pipeline；这里只负责编排多条流水线，
// 不引入 Manager，也不使用全局注册表。等出现真实业务 Processor 后，
// 在本函数中显式创建各自的 fileflow，并由 main 初始化流程统一触发。
func Fileflow() {
	cfg := global.GVA_CONFIG.Fileflow
	if !cfg.Enable {
		return
	}
	if _, err := BuildFileflowConfig(cfg); err != nil {
		zap.L().Error("failed to build fileflow config", zap.Error(err))
		return
	}

	// 当前项目尚未接入真实业务 Processor，所以下面仅保留未来接入示例。
	// 后续可按目录拆成 initXXXFileflow() 小函数，保持每条流水线边界清晰。
	//
	// ctx := context.Background()
	//
	// csvCfg, err := BuildFileflowConfig(global.GVA_CONFIG.Fileflow)
	// if err != nil {
	// 	zap.L().Error("failed to build csv fileflow config", zap.Error(err))
	// 	return
	// }
	// csvCfg.WatchDir = filepath.Clean("/data/input/csv")
	// csvCfg.Filter = fileflow.AndFilter(csvCfg.Filter, fileflow.ExtFilter(".csv"))
	// csvFlow, err := fileflow.New(csvCfg)
	// if err != nil {
	// 	zap.L().Error("failed to create csv fileflow", zap.Error(err))
	// 	return
	// }
	// csvFlow.Use(&CSVProcessor{})
	// StartFileflowPipeline(ctx, csvFlow)
	//
	// pcapCfg, err := BuildFileflowConfig(global.GVA_CONFIG.Fileflow)
	// if err != nil {
	// 	zap.L().Error("failed to build pcap fileflow config", zap.Error(err))
	// 	return
	// }
	// pcapCfg.WatchDir = filepath.Clean("/data/input/pcap")
	// pcapCfg.Filter = fileflow.AndFilter(pcapCfg.Filter, fileflow.ExtFilter(".pcap"))
	// pcapCfg.ReadyStrategy = fileflow.OKFileReady{}
	// pcapFlow, err := fileflow.New(pcapCfg)
	// if err != nil {
	// 	zap.L().Error("failed to create pcap fileflow", zap.Error(err))
	// 	return
	// }
	// pcapFlow.Use(&PcapArchiveProcessor{})
	// StartFileflowPipeline(ctx, pcapFlow)

	zap.L().Info("fileflow is enabled, waiting for explicit business pipelines")
}

// StopFileflow 保留兼容入口；v2 不再维护全局 fileflow 运行实例。
func StopFileflow() {}

// StartFileflowPipeline 使用传入的 ctx 独立启动一个 fileflow 流水线。
// 调用方负责为每条流水线维护自己的 ctx/cancel 生命周期。
func StartFileflowPipeline(ctx context.Context, p *fileflow.Pipeline) {
	if p == nil {
		return
	}
	go func() {
		err := p.Run(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			zap.L().Error("fileflow exited with error", zap.Error(err))
		}
	}()
}

// BuildFileflowConfig 将应用层 fileflow 配置转换为组件运行时配置。
func BuildFileflowConfig(cfg config.Fileflow) (fileflow.Config, error) {
	if strings.TrimSpace(cfg.WatchDir) == "" {
		return fileflow.Config{}, errors.New("fileflow: watch-dir is required when enabled")
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
	readyStrategy, err := parseReadyStrategy(cfg.ReadyStrategy, cfg.OKSuffix)
	if err != nil {
		return fileflow.Config{}, err
	}

	filter := func(name string) bool { return true }
	if cfg.IgnoreHidden {
		filter = func(name string) bool { return !strings.HasPrefix(name, ".") }
	}

	ffCfg := fileflow.Config{
		WatchDir:       filepath.Clean(cfg.WatchDir),
		Interval:       interval,
		Concurrency:    cfg.Concurrency,
		EventBuffer:    cfg.EventBuffer,
		ProcessTimeout: timeout,
		Filter:         filter,
		ReadyStrategy:  readyStrategy,
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

func parseDurationOrDefault(raw string, defaultValue time.Duration) (time.Duration, error) {
	if strings.TrimSpace(raw) == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(raw)
}

func parseReadyStrategy(raw, okSuffix string) (fileflow.ReadyStrategy, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "rename", "rename-ready":
		return fileflow.RenameReady{}, nil
	case "ok", "ok-file", "ok-file-ready":
		return fileflow.OKFileReady{Suffix: okSuffix}, nil
	default:
		return nil, fmt.Errorf("fileflow: unsupported ready-strategy %q", raw)
	}
}

func parseAfterProcess(raw string) (fileflow.AfterWriteAction, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "move", "move-to-processed", "archive", "archive-only":
		return fileflow.MoveToProcessed, nil
	case "delete":
		return fileflow.DeleteSource, nil
	case "keep":
		return fileflow.KeepSource, nil
	default:
		return 0, fmt.Errorf("fileflow: unsupported after-process action %q", raw)
	}
}
