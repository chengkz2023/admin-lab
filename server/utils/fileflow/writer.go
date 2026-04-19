package fileflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ResultWriter struct {
	outputDir     string
	afterWrite    AfterWriteAction
	doneDir       string
	postProcessor PostProcessor
}

// Write 将结果文件原子写入输出目录，并按配置处理源文件。
func (w *ResultWriter) Write(ctx context.Context, r Result) error {
	var written []string

	for _, f := range r.Files {
		if err := ctx.Err(); err != nil {
			return err
		}
		// 防止路径穿越和非法输出路径。
		if err := validateOutputPath(f); err != nil {
			return err
		}

		dir := w.outputDir
		if f.SubDir != "" {
			dir = filepath.Join(dir, filepath.Clean(f.SubDir))
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("fileflow: mkdir %s: %w", dir, err)
		}

		dest := filepath.Join(dir, f.Name)
		// 先写临时文件再 rename，避免下游读取半截文件。
		if err := writeAtomic(dest, f.Content); err != nil {
			return err
		}
		written = append(written, dest)
	}

	switch w.afterWrite {
	case MoveToProcessed:
		if err := os.MkdirAll(w.doneDir, 0o755); err == nil {
			_ = moveToDir(r.SourceEvent.Path, w.doneDir)
		}
	case DeleteSource:
		_ = os.Remove(r.SourceEvent.Path)
	}

	if w.postProcessor != nil {
		return w.postProcessor.PostProcess(ctx, written, r.Metadata)
	}
	return nil
}

// validateOutputPath 校验输出文件名和子目录，防止写出到 outputDir 之外。
func validateOutputPath(f OutputFile) error {
	if f.Name == "" {
		return errors.New("fileflow: output file name cannot be empty")
	}
	if f.Name != filepath.Base(f.Name) {
		return fmt.Errorf("fileflow: output file name %q must not contain path separator", f.Name)
	}
	if f.SubDir == "" {
		return nil
	}
	clean := filepath.Clean(f.SubDir)
	if clean == "." {
		return nil
	}
	if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
		return fmt.Errorf("fileflow: invalid subdir %q", f.SubDir)
	}
	return nil
}

// writeAtomic 在目标目录内创建临时文件并原子替换目标文件。
func writeAtomic(dest string, content []byte) error {
	dir := filepath.Dir(dest)
	tmpFile, err := os.CreateTemp(dir, "."+filepath.Base(dest)+".*.tmp")
	if err != nil {
		return fmt.Errorf("fileflow: create tmp for %s: %w", dest, err)
	}
	tmp := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmp)
	}()

	if _, err := tmpFile.Write(content); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("fileflow: write tmp %s: %w", tmp, err)
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

// moveToDir 将文件移动到指定目录，重名时自动追加时间戳。
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

// moveFile 优先 rename，跨文件系统时降级为 copy+remove。
func moveFile(src, dest string) error {
	if err := os.Rename(src, dest); err == nil {
		return nil
	}
	if err := copyFile(src, dest); err != nil {
		return err
	}
	return os.Remove(src)
}

// copyFile 复制单个文件并刷盘。
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
