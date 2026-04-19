# 目录文件处理流水线（工具化复用）

## 目标

把“持续处理目录文件”的通用能力抽成工具库，业务侧只关注：

1. 输入目录、输出目录等参数配置。
2. 业务处理逻辑（处理器）实现。
3. 在 API、定时任务、消息消费者中按同一方式调用。

## 工具包位置

- `server/utils/filepipeline/types.go`
- `server/utils/filepipeline/engine.go`
- `server/utils/filepipeline/processors.go`

## 核心能力

- 目录扫描（支持 `filePattern`）。
- 文件稳定性检查（防止读到半写入文件）。
- 文件认领（重命名为 `.processing` 防止重复处理）。
- 处理结果输出到 `outputDir`。
- 成功归档到 `archiveDir`。
- 失败转移到 `errorDir` 并记录失败明细。

## 开箱即用处理器

- `copy`：复制文件到输出目录（支持输出后缀）。
- `uppercase-text`：文本文件转大写后输出，非文本回退为复制。

## 像工具库一样使用

```go
package yourservice

import (
    "context"

    "github.com/flipped-aurora/gin-vue-admin/server/utils/filepipeline"
)

var pipelineEngine = filepipeline.NewDefaultEngine()

func runPipeline(ctx context.Context) error {
    _, err := pipelineEngine.RunOnce(ctx, filepipeline.Config{
        InputDir:     "D:/runtime/in",
        OutputDir:    "D:/runtime/out",
        FilePattern:  "*.txt",
        MaxFiles:     20,
        StableWaitMS: 1000,
        Processor:    "copy",
        OutputSuffix: "_processed",
    })
    return err
}
```

## 自定义处理器

1. 实现 `filepipeline.Processor` 接口。
2. 在引擎初始化时注册：`engine.RegisterProcessor(yourProcessor)`。
3. 调用 `RunOnce` 时传 `Processor: "your-key"`。

## API 适配层（当前项目）

- `GET /dirPipeline/profile`
- `POST /dirPipeline/runOnce`

`server/service/lab/reusable/dir_file_pipeline.go` 现在仅做请求/响应与 `filepipeline` 结果映射，不再承载底层编排细节。
