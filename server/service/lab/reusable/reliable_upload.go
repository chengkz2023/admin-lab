package reusable

import reusableRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"

type ReliableUploadService struct{}

func (r *ReliableUploadService) GetProfile() reusableRes.ReliableUploadProfile {
	return reusableRes.ReliableUploadProfile{
		Title:          "Reliable Upload Kit",
		Classification: "复用组件",
		Summary:        "面向生产环境的后端文件上报框架，已作为独立开源项目维护。覆盖分钟任务、大任务和业务触发任务三种模型，以备份文件作为重试唯一数据源，接口与框架完全解耦，可按项目自行实现 Repo 层接入。",
		GithubURL:      "https://github.com/chengkz2023/reliable-upload-kit",
		PackagePath:    "github.com/chengkz2023/reliable-upload-kit",
		DocumentPath:   "docs/reusable-reliable-upload.md",
		CapabilityPoints: []string{
			"生产与上报解耦，可分别调度 producer 和 uploader。",
			"按 task_code 隔离并发，适合多个上报任务并行共存。",
			"备份文件是重试唯一数据源，避免直接依赖原始业务表重查。",
			"支持启动恢复、分钟补窗、大任务续跑和业务任务断点续传。",
			"通过 Registry 挂接 DataSource、Reporter 和 FileNamer，业务接入点清晰。",
		},
		TaskModels: []reusableRes.ReliableUploadTaskModel{
			{Key: "minute", Name: "分钟任务", Scene: "适合高频、小批量、按时间窗稳定扫描的上报场景。", TriggerMode: "定时调度", Highlights: []string{"自动按 IntervalMinutes 和 DelaySeconds 计算窗口。", "OnStartup 可补齐停机期间缺失窗口。", "每个窗口可拆成多个 chunk 生成多份文件。"}},
			{Key: "big", Name: "自定义大任务", Scene: "适合历史补数、批量补传、任意时间窗重跑。", TriggerMode: "业务显式触发", Highlights: []string{"用 task_code + window_start + window_end 做幂等实例。", "支持先生产后上传，也支持中断后续跑。", "保留批次级 record_count 方便对账。"}},
			{Key: "biz", Name: "业务触发任务", Scene: "适合审批单、结算批次、活动单据等按业务键触发的上报。", TriggerMode: "传入 trigger_key / trigger_payload", Highlights: []string{"建议对 task_code + trigger_key 建唯一键保证幂等。", "DataSource 可从上下文读取 BizTrigger。", "文件命名可拼接 trigger_key 等轻量业务信息。"}},
		},
		EngineEntries: []reusableRes.ReliableUploadEntry{
			{Name: "RunProducer", Description: "扫描所有启用的分钟任务并生成 pending 文件。"},
			{Name: "RunUploader", Description: "兼容入口，统一上传分钟、大任务、业务任务的 pending 文件。"},
			{Name: "OnStartup", Description: "分钟补窗 + 恢复 running 的大任务和业务任务。"},
			{Name: "RunBigTask", Description: "按自定义时间窗口创建或恢复一个大任务实例。"},
			{Name: "RunBizTask", Description: "按 trigger_key 创建或恢复一个业务触发任务实例。"},
			{Name: "ProduceForTask", Description: "为指定 task_code 显式生产某个时间窗口的数据文件。"},
		},
		Interfaces: []reusableRes.ReliableUploadInterface{
			{Name: "DataSource", Signature: "CountChunks(ctx, cfg, start, end) / FetchChunk(ctx, cfg, start, end, index)", Notes: []string{"只负责稳定产出 chunk 总数和按索引取 chunk。", "业务任务可通过 BizTriggerFromContext 读取触发参数。"}},
			{Name: "Reporter", Signature: "Upload(ctx, cfg, item)", Notes: []string{"只处理文件上报，不关心 chunk 如何生成。", "UploadItem 内含 fileName、data、bizKey、meta 和 backupPath。"}},
			{Name: "TaskConfigRepo / UploadLogRepo / BigTaskRepo / BizTaskRepo", Signature: "配置与状态持久化仓储接口", Notes: []string{"框架不绑定具体 ORM，可按项目使用 GORM/MySQL 实现。", "admin-lab 当前先收录框架与接入说明，不强行引入具体表结构到系统库。"}},
		},
		IntegrationSteps: []string{
			"go get github.com/chengkz2023/reliable-upload-kit 引入依赖。",
			"在业务项目中实现 TaskConfigRepo 和对应状态表仓储。",
			"为每个 task_code 注册 DataSource、Reporter，必要时补 FileNamer。",
			"初始化 Engine，并在调度器中分别调用 RunProducer / RunUploader / OnStartup。",
			"对大任务和业务任务提供显式触发入口，例如管理台按钮或业务事件钩子。",
			"把 backup 目录纳入运维监控，确保备份文件生命周期与重试策略匹配。",
		},
		MigrationNotes: []string{
			"在 go.mod 中引入 github.com/chengkz2023/reliable-upload-kit，再按实际库表实现 Repo 层。",
			"如果内网已有上传/调度基础设施，可只复用 Engine、Registry、types 和接口契约。",
			"backup 路径、远端目录、重试上限、文件命名规则建议全部改为项目配置项。",
			"业务触发任务务必补唯一键和审计日志，避免重复触发带来脏数据。",
		},
		ExampleConfigs: []reusableRes.ReliableUploadExampleConfig{
			{TaskCode: "order_minute", TaskType: "minute", IntervalMinutes: 5, DelaySeconds: 60, BatchSize: 500, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order"},
			{TaskCode: "order_big", TaskType: "big", BatchSize: 2000, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order_big"},
			{TaskCode: "order_biz", TaskType: "biz", BatchSize: 2000, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order_biz"},
		},
		IncludedFiles: []reusableRes.ReliableUploadIncludedFile{
			{Path: "types.go", Role: "任务类型、状态、实体定义"},
			{Path: "interfaces.go", Role: "DataSource、Reporter、Repo 等核心契约"},
			{Path: "registry.go", Role: "task_code 级别注册中心"},
			{Path: "engine.go", Role: "生产、上传、恢复、重试主流程"},
			{Path: "biz_context.go", Role: "业务触发上下文透传"},
			{Path: "backup_fs.go", Role: "本地文件备份存储实现"},
		},
		CodeSnippets: []reusableRes.ReliableUploadCodeSnippet{
			{Title: "核心接口", Language: "go", Code: "type DataSource interface {\n    CountChunks(ctx context.Context, cfg TaskConfig, start, end time.Time) (int, error)\n    FetchChunk(ctx context.Context, cfg TaskConfig, start, end time.Time, index int) (Chunk, error)\n}\n\ntype Reporter interface {\n    Upload(ctx context.Context, cfg TaskConfig, item UploadItem) error\n}"},
			{Title: "引擎初始化", Language: "go", Code: "import kit \"github.com/chengkz2023/reliable-upload-kit\"\n\nregistry := kit.NewRegistry()\nregistry.RegisterDataSource(\"order_minute\", ds)\nregistry.RegisterReporter(\"order_minute\", rp)\n\nengine := kit.NewEngine(\n    registry,\n    cfgRepo,\n    uploadLogRepo,\n    bigRepo,\n    bizRepo,\n    kit.NewFSBackupStore(\"./backup\"),\n)\n\n_ = engine.RunProducer(ctx)\n_ = engine.RunUploader(ctx)"},
			{Title: "业务触发任务", Language: "go", Code: "trigger, ok := kit.BizTriggerFromContext(ctx)\nif ok {\n    fmt.Println(trigger.Key, trigger.Payload)\n}\n\n_ = engine.RunBizTask(ctx, \"order_biz\", \"approval_1001\", `{\"operator\":\"ops\"}`)"},
		},
	}
}
