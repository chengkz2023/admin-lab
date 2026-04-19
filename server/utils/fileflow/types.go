package fileflow

import (
	"context"
	"time"
)

// FileEvent 表示 watcher 识别出的“可处理文件事件”。
type FileEvent struct {
	Path      string
	Name      string
	Size      int64
	ModTime   time.Time
	CreatedAt time.Time
}

// OutputFile 表示 Processor 产出的单个目标文件。
type OutputFile struct {
	Name    string
	Content []byte
	SubDir  string
}

// Result 是 Processor.Process 的返回结构。
type Result struct {
	Files       []OutputFile
	Metadata    map[string]any
	SourceEvent FileEvent
}

// Processor 是业务方唯一需要实现的核心接口。
type Processor interface {
	Match(event FileEvent) bool
	Process(ctx context.Context, event FileEvent) (Result, error)
}

// PostProcessor 是结果文件写盘后的可选后置处理钩子。
type PostProcessor interface {
	PostProcess(ctx context.Context, written []string, meta map[string]any) error
}
