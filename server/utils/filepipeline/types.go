package filepipeline

import (
	"context"
	"sort"
	"strings"
)

const (
	DefaultPattern      = "*"
	DefaultMaxFiles     = 20
	DefaultStableWaitMS = 1000
	DefaultProcessor    = "copy"
	DefaultOutputSuffix = "_processed"
)

type Config struct {
	InputDir     string
	OutputDir    string
	ErrorDir     string
	ArchiveDir   string
	FilePattern  string
	MaxFiles     int
	StableWaitMS int
	Processor    string
	OutputSuffix string
}

type Result struct {
	EffectiveConfig Config
	Scanned         int
	Skipped         int
	Processed       int
	Failed          int
	OutputFiles     []string
	ArchivedFiles   []string
	SkippedFiles    []string
	FailureItems    []FailureItem
}

type FailureItem struct {
	File      string
	Stage     string
	Reason    string
	ErrorFile string
}

type Input struct {
	OriginalPath string
	OriginalName string
	ClaimedPath  string
}

type Output struct {
	SourcePath string
	OutputName string
}

type ProcessorInfo struct {
	Key         string
	Name        string
	Description string
}

type Processor interface {
	Info() ProcessorInfo
	Process(ctx context.Context, input Input, workDir string, cfg Config) ([]Output, error)
}

type Engine struct {
	processors map[string]Processor
}

func NewEngine(processors ...Processor) *Engine {
	engine := &Engine{
		processors: make(map[string]Processor),
	}
	for _, processor := range processors {
		_ = engine.RegisterProcessor(processor)
	}
	return engine
}

func NewDefaultEngine() *Engine {
	return NewEngine(
		NewCopyProcessor(),
		NewUppercaseTextProcessor(),
	)
}

func (e *Engine) RegisterProcessor(processor Processor) error {
	if processor == nil {
		return nil
	}
	info := processor.Info()
	key := strings.ToLower(strings.TrimSpace(info.Key))
	if key == "" {
		return nil
	}
	e.processors[key] = processor
	return nil
}

func (e *Engine) ProcessorInfos() []ProcessorInfo {
	items := make([]ProcessorInfo, 0, len(e.processors))
	for _, processor := range e.processors {
		items = append(items, processor.Info())
	}
	sort.Slice(items, func(i int, j int) bool {
		return items[i].Key < items[j].Key
	})
	return items
}
