# Reliable Upload Framework

## 模块归类

- 分区：复用组件
- 类型：后端可靠文件上报框架
- 收录位置：`server/utils/reliableupload`

## 背景与目标

`smart-upload` 的核心价值不在单个业务接口，而在一套可重复接入的可靠上报骨架：

- 数据生产与文件上报解耦
- 失败重试基于备份文件而不是反查业务表
- 支持分钟任务、大任务、业务触发任务三种模型
- 支持启动恢复、补窗、断点续跑

本次收录到 `admin-lab` 的目标，是先把这套骨架作为“可迁移资产”沉淀下来，方便后续搬到内网项目或继续扩展成完整演示。

## 当前收录范围

已收录文件：

- `server/utils/reliableupload/types.go`
- `server/utils/reliableupload/interfaces.go`
- `server/utils/reliableupload/registry.go`
- `server/utils/reliableupload/engine.go`
- `server/utils/reliableupload/biz_context.go`
- `server/utils/reliableupload/backup_fs.go`

配套展示入口：

- 后端资料接口：`/reliableUpload/profile`
- 前端页面：`实验室 / 复用组件 / 可靠上报框架`

## 三种任务模型

### 1. 分钟任务

适合高频、小批量、按固定时间窗口扫描的上报任务。

典型特点：

- 自动按 `IntervalMinutes` 和 `DelaySeconds` 计算窗口
- 支持停机后补齐缺失窗口
- 适合对账、流水、订单增量等分钟级上报

### 2. 自定义大任务

适合历史补数、批量补传、任意时间范围重跑。

典型特点：

- 使用 `task_code + window_start + window_end` 做实例幂等
- 支持大窗口拆批
- 支持生成后中断、恢复后继续上传

### 3. 业务触发任务

适合审批单、结算批次、活动单据等按业务键显式触发的上报。

典型特点：

- 使用 `trigger_key` 做实例幂等
- `DataSource` 可从上下文读取触发参数
- 适合手工触发、业务事件驱动、补偿式上报

## 核心接口

最小接入面主要有两类：

### DataSource

负责定义“怎么产出数据块”：

```go
type DataSource interface {
    CountChunks(ctx context.Context, cfg TaskConfig, start, end time.Time) (int, error)
    FetchChunk(ctx context.Context, cfg TaskConfig, start, end time.Time, index int) (Chunk, error)
}
```

### Reporter

负责定义“怎么把文件上报出去”：

```go
type Reporter interface {
    Upload(ctx context.Context, cfg TaskConfig, item UploadItem) error
}
```

除此之外，项目侧还需要实现：

- `TaskConfigRepo`
- `UploadLogRepo`
- `BigTaskRepo`
- `BizTaskRepo`

这些接口负责配置读取与状态持久化，框架本身不强绑具体 ORM 和表结构。

## 引擎入口

常用入口如下：

- `RunProducer`
- `RunUploader`
- `RunMinuteUploader`
- `RunBigUploader`
- `RunBizUploader`
- `OnStartup`
- `RunBigTask`
- `RunBizTask`
- `ProduceForTask`

推荐调度方式：

1. 系统启动时调用 `OnStartup`
2. 定时任务 A 调用 `RunProducer`
3. 定时任务 B 调用 `RunUploader`
4. 管理台或业务事件按需调用 `RunBigTask` / `RunBizTask`

## 在内网项目中的接入建议

建议按下面顺序迁入：

1. 先迁移 `server/utils/reliableupload`
2. 按内网项目规范实现 repo 和状态表
3. 为每个 `task_code` 注册 `DataSource` 与 `Reporter`
4. 接入调度系统
5. 再补管理台触发入口和监控告警

## 建议同步的表与约束

如果要在内网完整跑通，通常需要补这些持久化对象：

- 分钟任务日志表
- 大任务实例表
- 大任务批次表
- 业务任务实例表
- 业务任务批次表

建议重点保证这些约束：

- 分钟任务时间窗唯一性
- 大任务 `(task_code, window_start, window_end)` 唯一性
- 业务任务 `(task_code, trigger_key)` 唯一性

## 迁移时要特别注意

- `backup` 路径要改成项目配置，不能写死
- 文件命名规则要和内网归档规范对齐
- 远端目录、重试上限、批次大小建议都改成可配置项
- 业务触发任务要补审计日志，避免重复触发难追踪
- 如果已有调度或上传组件，可以只复用这套引擎和接口契约

## 当前边界

这次在 `admin-lab` 中先完成的是：

- 框架代码收录
- 资料接口与展示页
- 菜单 / API / 权限种子
- 迁移文档

还没有做的是：

- 把 `smart-upload` 示例里的 MySQL repo 直接并入当前系统库
- 提供一套可直接跑任务的专用 demo 数据表
- 在页面上做真实触发和运行状态回显

如果后续要继续扩展，建议下一步补一个“可触发 demo 版 repo + 运行面板”。
