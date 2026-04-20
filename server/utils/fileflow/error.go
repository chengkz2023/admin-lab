package fileflow

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type ErrorHandler interface {
	Handle(ctx context.Context, event FileEvent, err error)
}

// DefaultErrorHandler 将失败源文件移到 failed 目录，并写 .error 诊断文件。
type DefaultErrorHandler struct {
	FailedDir string
	Logger    *slog.Logger
}

// Handle 执行默认失败处置逻辑。
func (h *DefaultErrorHandler) Handle(_ context.Context, e FileEvent, err error) {
	if h.Logger == nil {
		h.Logger = slog.Default()
	}
	h.Logger.Error("fileflow: process failed", "file", e.Name, "err", err)

	if mkErr := os.MkdirAll(h.FailedDir, 0o755); mkErr != nil {
		return
	}

	fileName := e.Name
	if fileName == "" {
		fileName = filepath.Base(e.Path)
	}
	dest := filepath.Join(h.FailedDir, fileName)
	if _, statErr := os.Stat(dest); statErr == nil {
		dest = filepath.Join(h.FailedDir, fmt.Sprintf("%s.%d", fileName, time.Now().UnixNano()))
	}
	_ = moveFile(e.Path, dest)

	errContent := fmt.Sprintf(
		"time:  %s\nfile:  %s\nerror: %v\n",
		time.Now().Format(time.RFC3339),
		e.Path,
		err,
	)
	_ = os.WriteFile(dest+".error", []byte(errContent), 0o644)
}
