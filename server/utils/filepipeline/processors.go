package filepipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type copyProcessor struct{}
type uppercaseTextProcessor struct{}

func NewCopyProcessor() Processor {
	return copyProcessor{}
}

func NewUppercaseTextProcessor() Processor {
	return uppercaseTextProcessor{}
}

func (p copyProcessor) Info() ProcessorInfo {
	return ProcessorInfo{
		Key:         "copy",
		Name:        "复制输出",
		Description: "将认领后的文件复制到输出目录，可配置输出后缀。",
	}
}

func (p copyProcessor) Process(_ context.Context, input Input, workDir string, cfg Config) ([]Output, error) {
	outputName := applyOutputSuffix(input.OriginalName, cfg.OutputSuffix)
	target := filepath.Join(workDir, outputName)
	if err := copyFile(input.ClaimedPath, target); err != nil {
		return nil, err
	}
	return []Output{
		{
			SourcePath: target,
			OutputName: outputName,
		},
	}, nil
}

func (p uppercaseTextProcessor) Info() ProcessorInfo {
	return ProcessorInfo{
		Key:         "uppercase-text",
		Name:        "文本转大写",
		Description: "将文本类文件内容转换为大写后输出，非文本文件回退为复制输出。",
	}
}

func (p uppercaseTextProcessor) Process(ctx context.Context, input Input, workDir string, cfg Config) ([]Output, error) {
	ext := strings.ToLower(filepath.Ext(input.OriginalName))
	if !isTextLike(ext) {
		return copyProcessor{}.Process(ctx, input, workDir, cfg)
	}

	data, err := os.ReadFile(input.ClaimedPath)
	if err != nil {
		return nil, errors.Wrap(err, "读取输入文件失败")
	}
	outputName := applyOutputSuffix(input.OriginalName, cfg.OutputSuffix)
	target := filepath.Join(workDir, outputName)
	if err = os.WriteFile(target, []byte(strings.ToUpper(string(data))), 0o644); err != nil {
		return nil, errors.Wrap(err, "写入处理结果文件失败")
	}
	return []Output{
		{
			SourcePath: target,
			OutputName: outputName,
		},
	}, nil
}

func isTextLike(ext string) bool {
	switch ext {
	case ".txt", ".log", ".csv", ".json", ".xml", ".md":
		return true
	default:
		return false
	}
}

func applyOutputSuffix(fileName string, suffix string) string {
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ext := filepath.Ext(fileName)
	return base + suffix + ext
}
