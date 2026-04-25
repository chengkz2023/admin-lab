package config

// Fileflow 是 fileflow 工具的配置结构。
type Fileflow struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	WatchDir string `mapstructure:"watch-dir" json:"watch-dir" yaml:"watch-dir"`
	// OutputDir 仅用于兼容旧配置加载。在 fileflow v2 中，
	// 输出由业务 Processor 自行写入，框架不再使用该字段。
	OutputDir     string `mapstructure:"output-dir" json:"output-dir" yaml:"output-dir"`
	FailedDir     string `mapstructure:"failed-dir" json:"failed-dir" yaml:"failed-dir"`
	ReadyStrategy string `mapstructure:"ready-strategy" json:"ready-strategy" yaml:"ready-strategy"`
	OKSuffix      string `mapstructure:"ok-suffix" json:"ok-suffix" yaml:"ok-suffix"`
	Interval      string `mapstructure:"interval" json:"interval" yaml:"interval"`
	Timeout       string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	BackoffBase   string `mapstructure:"backoff-base" json:"backoff-base" yaml:"backoff-base"`
	AfterProcess  string `mapstructure:"after-process" json:"after-process" yaml:"after-process"`
	Concurrency   int    `mapstructure:"concurrency" json:"concurrency" yaml:"concurrency"`
	MaxRetries    int    `mapstructure:"max-retries" json:"max-retries" yaml:"max-retries"`
	EventBuffer   int    `mapstructure:"event-buffer" json:"event-buffer" yaml:"event-buffer"`
	IgnoreHidden  bool   `mapstructure:"ignore-hidden" json:"ignore-hidden" yaml:"ignore-hidden"`
}
