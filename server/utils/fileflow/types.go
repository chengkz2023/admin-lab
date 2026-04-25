package fileflow

import (
	"context"
	"time"
)

// FileEvent 表示 watcher 发现的可处理源文件事件。
type FileEvent struct {
	Path      string
	Name      string
	Size      int64
	ModTime   time.Time
	CreatedAt time.Time
}

// Result 承载 Processor 传递给 Hook 的可选业务元数据。
// 框架不再写输出文件，输出由业务 Processor 自行负责。
type Result struct {
	Metadata map[string]any
}

// Processor 是业务方处理就绪源文件的核心接口。
type Processor interface {
	Process(ctx context.Context, event FileEvent) (Result, error)
}

// Hook 观察处理结果，并执行源文件清理或业务收尾动作。
type Hook interface {
	OnSuccess(ctx context.Context, event FileEvent, result Result) error
	OnError(ctx context.Context, event FileEvent, err error)
}

// HookFunc 将函数适配为 Hook，nil 回调会被忽略。
type HookFunc struct {
	Success func(ctx context.Context, event FileEvent, result Result) error
	Error   func(ctx context.Context, event FileEvent, err error)
}

func (h HookFunc) OnSuccess(ctx context.Context, event FileEvent, result Result) error {
	if h.Success == nil {
		return nil
	}
	return h.Success(ctx, event, result)
}

func (h HookFunc) OnError(ctx context.Context, event FileEvent, err error) {
	if h.Error != nil {
		h.Error(ctx, event, err)
	}
}
