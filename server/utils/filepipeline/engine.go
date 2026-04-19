package filepipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func NormalizeConfig(raw Config) Config {
	cfg := Config{
		InputDir:     cleanPath(raw.InputDir),
		OutputDir:    cleanPath(raw.OutputDir),
		ErrorDir:     cleanPath(raw.ErrorDir),
		ArchiveDir:   cleanPath(raw.ArchiveDir),
		FilePattern:  strings.TrimSpace(raw.FilePattern),
		MaxFiles:     raw.MaxFiles,
		StableWaitMS: raw.StableWaitMS,
		Processor:    strings.ToLower(strings.TrimSpace(raw.Processor)),
		OutputSuffix: strings.TrimSpace(raw.OutputSuffix),
	}

	if cfg.FilePattern == "" {
		cfg.FilePattern = DefaultPattern
	}
	if cfg.MaxFiles <= 0 {
		cfg.MaxFiles = DefaultMaxFiles
	}
	if cfg.StableWaitMS <= 0 {
		cfg.StableWaitMS = DefaultStableWaitMS
	}
	if cfg.Processor == "" {
		cfg.Processor = DefaultProcessor
	}
	if cfg.OutputSuffix == "" {
		cfg.OutputSuffix = DefaultOutputSuffix
	}
	if cfg.ErrorDir == "" && cfg.OutputDir != "" {
		cfg.ErrorDir = filepath.Join(cfg.OutputDir, "_errors")
	}
	if cfg.ArchiveDir == "" && cfg.InputDir != "" {
		cfg.ArchiveDir = filepath.Join(cfg.InputDir, "_archive")
	}
	return cfg
}

func (e *Engine) RunOnce(ctx context.Context, raw Config) (Result, error) {
	cfg := NormalizeConfig(raw)
	result := Result{
		EffectiveConfig: cfg,
		OutputFiles:     make([]string, 0),
		ArchivedFiles:   make([]string, 0),
		SkippedFiles:    make([]string, 0),
		FailureItems:    make([]FailureItem, 0),
	}

	if err := validateConfig(cfg); err != nil {
		return result, err
	}
	processor, ok := e.processors[cfg.Processor]
	if !ok {
		return result, errors.New("不支持的处理器: " + cfg.Processor)
	}

	for _, targetDir := range []string{cfg.OutputDir, cfg.ErrorDir, cfg.ArchiveDir} {
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return result, errors.Wrapf(err, "创建目录失败: %s", targetDir)
		}
	}

	candidates, err := collectCandidates(cfg.InputDir, cfg.FilePattern)
	if err != nil {
		return result, err
	}
	result.Scanned = len(candidates)
	if len(candidates) == 0 {
		return result, nil
	}

	limit := cfg.MaxFiles
	if limit > len(candidates) {
		limit = len(candidates)
	}
	candidates = candidates[:limit]

	for _, filePath := range candidates {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		stable, stableErr := waitUntilStable(filePath, time.Duration(cfg.StableWaitMS)*time.Millisecond)
		if stableErr != nil {
			result.Failed++
			result.FailureItems = append(result.FailureItems, FailureItem{
				File:   filePath,
				Stage:  "稳定性检查",
				Reason: stableErr.Error(),
			})
			continue
		}
		if !stable {
			result.Skipped++
			result.SkippedFiles = append(result.SkippedFiles, filePath)
			continue
		}

		claimedPath, claimErr := claimFile(filePath)
		if claimErr != nil {
			result.Failed++
			result.FailureItems = append(result.FailureItems, FailureItem{
				File:   filePath,
				Stage:  "认领文件",
				Reason: claimErr.Error(),
			})
			continue
		}

		input := Input{
			OriginalPath: filePath,
			OriginalName: filepath.Base(filePath),
			ClaimedPath:  claimedPath,
		}

		workDir, workErr := os.MkdirTemp("", "filepipeline-*")
		if workErr != nil {
			captureFailure(&result, cfg, input, "创建工作目录", workErr)
			continue
		}

		outputs, processErr := processor.Process(ctx, input, workDir, cfg)
		if processErr != nil {
			_ = os.RemoveAll(workDir)
			captureFailure(&result, cfg, input, "业务处理", processErr)
			continue
		}
		if len(outputs) == 0 {
			_ = os.RemoveAll(workDir)
			captureFailure(&result, cfg, input, "业务处理", errors.New("处理器未返回任何输出文件"))
			continue
		}

		allOutputDone := true
		for _, output := range outputs {
			if output.SourcePath == "" {
				allOutputDone = false
				captureFailure(&result, cfg, input, "输出文件", errors.New("输出源文件路径为空"))
				break
			}

			outputName := strings.TrimSpace(output.OutputName)
			if outputName == "" {
				outputName = filepath.Base(output.SourcePath)
			}

			targetPath, targetErr := uniqueTargetPath(cfg.OutputDir, outputName)
			if targetErr != nil {
				allOutputDone = false
				captureFailure(&result, cfg, input, "输出文件", targetErr)
				break
			}
			if copyErr := copyFile(output.SourcePath, targetPath); copyErr != nil {
				allOutputDone = false
				captureFailure(&result, cfg, input, "输出文件", copyErr)
				break
			}
			result.OutputFiles = append(result.OutputFiles, targetPath)
		}
		_ = os.RemoveAll(workDir)
		if !allOutputDone {
			continue
		}

		archivePath, archiveErr := uniqueTargetPath(cfg.ArchiveDir, input.OriginalName)
		if archiveErr != nil {
			captureFailure(&result, cfg, input, "归档文件", archiveErr)
			continue
		}
		if moveErr := moveFile(claimedPath, archivePath); moveErr != nil {
			captureFailure(&result, cfg, input, "归档文件", moveErr)
			continue
		}

		result.ArchivedFiles = append(result.ArchivedFiles, archivePath)
		result.Processed++
	}

	return result, nil
}

func captureFailure(result *Result, cfg Config, input Input, stage string, err error) {
	result.Failed++
	failure := FailureItem{
		File:   input.OriginalPath,
		Stage:  stage,
		Reason: err.Error(),
	}
	if input.ClaimedPath != "" {
		errorPath, moveErr := moveToErrorDir(input.ClaimedPath, cfg.ErrorDir, input.OriginalName)
		if moveErr == nil {
			failure.ErrorFile = errorPath
		} else {
			failure.Reason = failure.Reason + "；转移到错误目录失败: " + moveErr.Error()
		}
	}
	result.FailureItems = append(result.FailureItems, failure)
}

func validateConfig(cfg Config) error {
	if cfg.InputDir == "" {
		return errors.New("inputDir 不能为空")
	}
	if cfg.OutputDir == "" {
		return errors.New("outputDir 不能为空")
	}
	if info, err := os.Stat(cfg.InputDir); err != nil {
		return errors.Wrap(err, "inputDir 无法访问")
	} else if !info.IsDir() {
		return errors.New("inputDir 必须是目录")
	}

	if isSamePath(cfg.InputDir, cfg.OutputDir) {
		return errors.New("inputDir 与 outputDir 不能是同一目录")
	}
	if isSamePath(cfg.InputDir, cfg.ErrorDir) {
		return errors.New("inputDir 与 errorDir 不能是同一目录")
	}
	if isSamePath(cfg.InputDir, cfg.ArchiveDir) {
		return errors.New("inputDir 与 archiveDir 不能是同一目录")
	}
	if _, err := filepath.Match(cfg.FilePattern, "sample.txt"); err != nil {
		return errors.Wrap(err, "filePattern 非法")
	}
	return nil
}

func cleanPath(raw string) string {
	path := strings.TrimSpace(raw)
	if path == "" {
		return ""
	}
	return filepath.Clean(path)
}

func isSamePath(left string, right string) bool {
	return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
}

func collectCandidates(inputDir string, pattern string) ([]string, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, errors.Wrap(err, "读取 inputDir 失败")
	}
	results := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".processing") {
			continue
		}
		matched, matchErr := filepath.Match(pattern, name)
		if matchErr != nil {
			return nil, errors.Wrap(matchErr, "匹配 filePattern 失败")
		}
		if !matched {
			continue
		}
		results = append(results, filepath.Join(inputDir, name))
	}
	sort.Strings(results)
	return results, nil
}

func waitUntilStable(path string, wait time.Duration) (bool, error) {
	before, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "等待前获取文件状态失败")
	}
	time.Sleep(wait)
	after, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "等待后获取文件状态失败")
	}
	return before.Size() == after.Size() && before.ModTime().Equal(after.ModTime()), nil
}

func claimFile(path string) (string, error) {
	claimedPath := fmt.Sprintf("%s.%d.processing", path, time.Now().UnixNano())
	if err := os.Rename(path, claimedPath); err != nil {
		return "", errors.Wrap(err, "认领文件时重命名失败")
	}
	return claimedPath, nil
}

func moveToErrorDir(claimedPath string, errorDir string, originalName string) (string, error) {
	targetName := addTimestampSuffix(originalName, "_failed")
	targetPath, err := uniqueTargetPath(errorDir, targetName)
	if err != nil {
		return "", err
	}
	if moveErr := moveFile(claimedPath, targetPath); moveErr != nil {
		return "", moveErr
	}
	return targetPath, nil
}

func addTimestampSuffix(fileName string, suffix string) string {
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ext := filepath.Ext(fileName)
	return fmt.Sprintf("%s%s_%d%s", base, suffix, time.Now().UnixNano(), ext)
}

func uniqueTargetPath(dir string, fileName string) (string, error) {
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ext := filepath.Ext(fileName)
	for index := 0; index < 10000; index++ {
		candidate := fileName
		if index > 0 {
			candidate = fmt.Sprintf("%s_%d%s", base, index, ext)
		}
		fullPath := filepath.Join(dir, candidate)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath, nil
		} else if err != nil {
			return "", errors.Wrap(err, "检查输出路径失败")
		}
	}
	return "", errors.New("无法分配唯一目标路径")
}

func moveFile(from string, to string) error {
	if err := os.Rename(from, to); err == nil {
		return nil
	} else if !isCrossDeviceRenameErr(err) {
		return errors.Wrap(err, "移动文件失败")
	}

	if err := copyFile(from, to); err != nil {
		return err
	}
	if err := os.Remove(from); err != nil {
		return errors.Wrap(err, "跨盘复制后删除源文件失败")
	}
	return nil
}

func isCrossDeviceRenameErr(err error) bool {
	linkErr, ok := err.(*os.LinkError)
	if !ok {
		return false
	}
	if errors.Is(linkErr.Err, syscall.EXDEV) {
		return true
	}
	return strings.Contains(strings.ToLower(linkErr.Err.Error()), "cross-device")
}

func copyFile(from string, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return errors.Wrap(err, "打开源文件失败")
	}
	defer source.Close()

	target, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return errors.Wrap(err, "打开目标文件失败")
	}
	defer target.Close()

	if _, err = io.Copy(target, source); err != nil {
		return errors.Wrap(err, "复制文件内容失败")
	}
	return nil
}
