package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BootstrapConfig 从 yaml 读取的基础配置（只用于启动阶段）
// 对应你的 yaml 文件中的数据库连接等基础字段
type BootstrapConfig struct {
	MySQL struct {
		DSN string `yaml:"dsn"` // "user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	} `yaml:"mysql"`
	App struct {
		UploadSource string        `yaml:"upload_source"` // 本节点唯一标志，如 "smms-node-01"
		PollInterval time.Duration `yaml:"poll_interval"` // 热更新间隔，如 "30s"
	} `yaml:"app"`
}

// Bootstrap 在 main.go 中调用：读取 yaml → 初始化 DB → 启动热更新管理器
//
//	cfg, err := loadYaml("config.yaml")  // 你已有的 yaml 加载逻辑
//	if err != nil { ... }
//	db, err := config.Bootstrap(cfg)
//	if err != nil { ... }
func Bootstrap(cfg *BootstrapConfig) (*gorm.DB, error) {
	if cfg.App.UploadSource == "" {
		return nil, fmt.Errorf("config: app.upload_source must not be empty")
	}
	if cfg.App.PollInterval <= 0 {
		cfg.App.PollInterval = 30 * time.Second
	}

	// 初始化 GORM（使用你项目已有的 DB 实例也可以，直接传入 InitManager 即可）
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("config: connect mysql failed: %w", err)
	}

	// 自动迁移（可选：若表已建好则删除此行）
	// _ = db.AutoMigrate(&UploadConfig{})

	// 启动热更新管理器
	if err := InitManager(db, cfg.App.UploadSource, cfg.App.PollInterval); err != nil {
		return nil, err
	}
	return db, nil
}
