package fileflow

import (
	"context"
	"os"
	"strings"
)

// ReadyStrategy 判断发现的文件是否完整且可处理。
type ReadyStrategy interface {
	Ready(ctx context.Context, path string, info os.FileInfo) (bool, error)
}

// RenameReady 认为以最终文件名出现的文件已就绪，适用于
// 上游先写临时文件、再原子 rename 的投递方式。
type RenameReady struct{}

func (RenameReady) Ready(_ context.Context, _ string, _ os.FileInfo) (bool, error) {
	return true, nil
}

// OKFileReady 仅在旁路信令文件存在时认为源文件已就绪。
// 例如 data.csv.ok 存在时，data.csv 才会被处理。
type OKFileReady struct {
	Suffix string
}

func (r OKFileReady) Ready(_ context.Context, path string, _ os.FileInfo) (bool, error) {
	suffix := r.Suffix
	if suffix == "" {
		suffix = ".ok"
	}
	if strings.HasSuffix(path, suffix) {
		return false, nil
	}
	_, err := os.Stat(path + suffix)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
