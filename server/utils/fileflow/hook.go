package fileflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// ChainHook 按顺序执行 Hook；成功流程遇到首个错误即停止，
// 错误流程会广播给每个 Hook。
func ChainHook(hooks ...Hook) Hook {
	return chainHook(hooks)
}

type chainHook []Hook

func (h chainHook) OnSuccess(ctx context.Context, event FileEvent, result Result) error {
	for _, hook := range h {
		if hook == nil {
			continue
		}
		if err := hook.OnSuccess(ctx, event, result); err != nil {
			return err
		}
	}
	return nil
}

func (h chainHook) OnError(ctx context.Context, event FileEvent, err error) {
	for _, hook := range h {
		if hook != nil {
			hook.OnError(ctx, event, err)
		}
	}
}

// DefaultHook 处理源文件归档、删除和失败隔离。
type DefaultHook struct {
	DoneDir    string
	FailedDir  string
	AfterWrite AfterWriteAction
	Logger     *slog.Logger
}

func (h *DefaultHook) OnSuccess(_ context.Context, event FileEvent, _ Result) error {
	switch h.AfterWrite {
	case MoveToProcessed, ArchiveOnly:
		if h.DoneDir == "" {
			return nil
		}
		if err := os.MkdirAll(h.DoneDir, 0o755); err != nil {
			return fmt.Errorf("fileflow: mkdir done dir %s: %w", h.DoneDir, err)
		}
		return moveToDir(event.Path, h.DoneDir)
	case DeleteSource:
		return os.Remove(event.Path)
	case KeepSource:
		return nil
	default:
		return nil
	}
}

func (h *DefaultHook) OnError(_ context.Context, event FileEvent, err error) {
	logger := h.Logger
	if logger == nil {
		logger = slog.Default()
	}
	logger.Error("fileflow: process failed", "file", event.Name, "err", err)

	if h.FailedDir == "" {
		return
	}
	if mkErr := os.MkdirAll(h.FailedDir, 0o755); mkErr != nil {
		logger.Error("fileflow: mkdir failed dir failed", "dir", h.FailedDir, "err", mkErr)
		return
	}

	fileName := event.Name
	if fileName == "" {
		fileName = filepath.Base(event.Path)
	}
	dest := filepath.Join(h.FailedDir, fileName)
	if _, statErr := os.Stat(dest); statErr == nil {
		dest = filepath.Join(h.FailedDir, fmt.Sprintf("%s.%d", fileName, time.Now().UnixNano()))
	}
	if moveErr := moveFile(event.Path, dest); moveErr != nil {
		logger.Error("fileflow: move failed file failed", "file", event.Path, "dest", dest, "err", moveErr)
		return
	}

	errContent := fmt.Sprintf(
		"time:  %s\nfile:  %s\nerror: %v\n",
		time.Now().Format(time.RFC3339),
		event.Path,
		err,
	)
	if writeErr := os.WriteFile(dest+".error", []byte(errContent), 0o644); writeErr != nil {
		logger.Error("fileflow: write error file failed", "file", dest+".error", "err", writeErr)
	}
}

// AtomicWriteFile 先写临时文件，再 rename 到目标路径。
func AtomicWriteFile(dest string, content []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("fileflow: mkdir %s: %w", filepath.Dir(dest), err)
	}
	tmpFile, err := os.CreateTemp(filepath.Dir(dest), "."+filepath.Base(dest)+".*.tmp")
	if err != nil {
		return fmt.Errorf("fileflow: create tmp for %s: %w", dest, err)
	}
	tmp := tmpFile.Name()
	defer func() { _ = os.Remove(tmp) }()

	if _, err := tmpFile.Write(content); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("fileflow: write tmp %s: %w", tmp, err)
	}
	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("fileflow: chmod tmp %s: %w", tmp, err)
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("fileflow: sync tmp %s: %w", tmp, err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("fileflow: close tmp %s: %w", tmp, err)
	}
	if err := os.Remove(dest); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("fileflow: remove old dest %s: %w", dest, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		return fmt.Errorf("fileflow: rename %s to %s: %w", tmp, dest, err)
	}
	return nil
}

func moveToDir(src, dir string) error {
	if src == "" {
		return nil
	}
	base := filepath.Base(src)
	dest := filepath.Join(dir, base)
	if _, err := os.Stat(dest); err == nil {
		dest = filepath.Join(dir, fmt.Sprintf("%s.%d", base, time.Now().UnixNano()))
	}
	return moveFile(src, dest)
}

func moveFile(src, dest string) error {
	if err := os.Rename(src, dest); err == nil {
		return nil
	}
	if err := copyFile(src, dest); err != nil {
		return err
	}
	return os.Remove(src)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
