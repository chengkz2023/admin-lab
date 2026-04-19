package reusable

import (
	"context"

	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	reusableRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/filepipeline"
)

type DirFilePipelineService struct{}

var dirFilePipelineEngine = filepipeline.NewDefaultEngine()

func (s *DirFilePipelineService) GetProfile() reusableRes.DirFilePipelineProfile {
	defaultCfg := filepipeline.NormalizeConfig(filepipeline.Config{
		InputDir:     "./runtime/in",
		OutputDir:    "./runtime/out",
		FilePattern:  "*.txt",
		MaxFiles:     filepipeline.DefaultMaxFiles,
		StableWaitMS: filepipeline.DefaultStableWaitMS,
		Processor:    filepipeline.DefaultProcessor,
		OutputSuffix: filepipeline.DefaultOutputSuffix,
	})

	return reusableRes.DirFilePipelineProfile{
		Title:          "目录文件处理流水线",
		Classification: "复用组件",
		Summary: "基于 utils/filepipeline 封装的可复用目录处理引擎。你只需配置目录参数并实现处理器，" +
			"即可像使用开源工具库一样复用整套扫描、认领、输出、归档与失败转移流程。",
		Highlights: []string{
			"引擎位于 server/utils/filepipeline，可在任意 service 或任务中复用。",
			"支持按 filePattern 扫描，并在处理前做稳定性检查。",
			"通过重命名认领文件，避免并发场景下重复处理。",
			"内置输出、归档、失败目录管理，减少重复样板代码。",
			"处理器接口可插拔，只需聚焦业务逻辑实现。",
		},
		QuickSteps: []string{
			"创建 filepipeline.Engine（可直接使用 NewDefaultEngine）。",
			"按业务实现并注册 Processor。",
			"传入 Config 执行 RunOnce。",
			"在 API、定时任务或消费者中复用同一套调用方式。",
		},
		Processors: buildProcessorOptions(),
		DefaultConfig: reusableRes.DirFilePipelineConfigSample{
			InputDir:     defaultCfg.InputDir,
			OutputDir:    defaultCfg.OutputDir,
			ErrorDir:     defaultCfg.ErrorDir,
			ArchiveDir:   defaultCfg.ArchiveDir,
			FilePattern:  defaultCfg.FilePattern,
			MaxFiles:     defaultCfg.MaxFiles,
			StableWaitMs: defaultCfg.StableWaitMS,
			Processor:    defaultCfg.Processor,
			OutputSuffix: defaultCfg.OutputSuffix,
		},
		IntegrationNotes: []string{
			"建议将引擎实例放在 service 包级变量，避免重复初始化。",
			"自定义处理器实现 filepipeline.Processor 接口并注册到引擎。",
			"业务校验放在处理器内，流程编排统一由引擎维护。",
			"可从 API、定时任务、队列消费者触发 RunOnce。",
			"建议处理逻辑保持幂等，便于失败重试。",
		},
	}
}

func (s *DirFilePipelineService) RunOnce(ctx context.Context, req reusableReq.DirFilePipelineRunRequest) (reusableRes.DirFilePipelineRunResult, error) {
	runResult, err := dirFilePipelineEngine.RunOnce(ctx, filepipeline.Config{
		InputDir:     req.InputDir,
		OutputDir:    req.OutputDir,
		ErrorDir:     req.ErrorDir,
		ArchiveDir:   req.ArchiveDir,
		FilePattern:  req.FilePattern,
		MaxFiles:     req.MaxFiles,
		StableWaitMS: req.StableWaitMs,
		Processor:    req.Processor,
		OutputSuffix: req.OutputSuffix,
	})
	return mapRunResult(runResult), err
}

func buildProcessorOptions() []reusableRes.DirFilePipelineProcessorOption {
	infos := dirFilePipelineEngine.ProcessorInfos()
	items := make([]reusableRes.DirFilePipelineProcessorOption, 0, len(infos))
	for _, info := range infos {
		items = append(items, reusableRes.DirFilePipelineProcessorOption{
			Key:         info.Key,
			Name:        info.Name,
			Description: info.Description,
		})
	}
	return items
}

func mapRunResult(from filepipeline.Result) reusableRes.DirFilePipelineRunResult {
	failures := make([]reusableRes.DirFilePipelineFailureItem, 0, len(from.FailureItems))
	for _, item := range from.FailureItems {
		failures = append(failures, reusableRes.DirFilePipelineFailureItem{
			File:      item.File,
			Stage:     item.Stage,
			Reason:    item.Reason,
			ErrorFile: item.ErrorFile,
		})
	}

	return reusableRes.DirFilePipelineRunResult{
		EffectiveConfig: reusableRes.DirFilePipelineConfigSample{
			InputDir:     from.EffectiveConfig.InputDir,
			OutputDir:    from.EffectiveConfig.OutputDir,
			ErrorDir:     from.EffectiveConfig.ErrorDir,
			ArchiveDir:   from.EffectiveConfig.ArchiveDir,
			FilePattern:  from.EffectiveConfig.FilePattern,
			MaxFiles:     from.EffectiveConfig.MaxFiles,
			StableWaitMs: from.EffectiveConfig.StableWaitMS,
			Processor:    from.EffectiveConfig.Processor,
			OutputSuffix: from.EffectiveConfig.OutputSuffix,
		},
		Scanned:       from.Scanned,
		Skipped:       from.Skipped,
		Processed:     from.Processed,
		Failed:        from.Failed,
		OutputFiles:   from.OutputFiles,
		ArchivedFiles: from.ArchivedFiles,
		SkippedFiles:  from.SkippedFiles,
		FailureItems:  failures,
	}
}
