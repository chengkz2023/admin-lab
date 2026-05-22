package config

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// UploadConfigMap 以 upload_source+code 为 key 的配置快照
type UploadConfigMap map[string]*UploadConfig

// configManager 热插拔配置管理器（内部单例）
type configManager struct {
	mu           sync.RWMutex
	db           *gorm.DB
	uploadSource string        // 当前服务节点的唯一标志（从 yaml 初始化时注入）
	configs      UploadConfigMap
	interval     time.Duration // 轮询间隔
	lastCheck    time.Time     // 上次检测时间
	onChange     []func(UploadConfigMap) // 变更回调钩子列表
	cancel       context.CancelFunc
}

var (
	globalManager *configManager
	once          sync.Once
)

// InitManager 在项目启动时调用，传入已初始化的 GORM DB、当前节点标识、轮询间隔
// uploadSource: 与表中 upload_source 字段对应，用于只拉取本节点配置
// pollInterval: 建议 10s ~ 60s
func InitManager(db *gorm.DB, uploadSource string, pollInterval time.Duration) error {
	var initErr error
	once.Do(func() {
		globalManager = &configManager{
			db:           db,
			uploadSource: uploadSource,
			configs:      make(UploadConfigMap),
			interval:     pollInterval,
		}
		// 启动时同步加载一次，确保配置立即可用
		if err := globalManager.load(); err != nil {
			initErr = fmt.Errorf("config: initial load failed: %w", err)
			return
		}
		// 启动后台热更新 goroutine
		ctx, cancel := context.WithCancel(context.Background())
		globalManager.cancel = cancel
		go globalManager.watch(ctx)
	})
	return initErr
}

// GetManager 获取全局管理器（InitManager 必须先调用）
func GetManager() *configManager {
	if globalManager == nil {
		panic("config: manager not initialized, call InitManager first")
	}
	return globalManager
}

// ─── 公开查询方法 ────────────────────────────────────────────────────────────

// All 返回当前所有配置的只读快照
func All() UploadConfigMap {
	m := GetManager()
	m.mu.RLock()
	defer m.mu.RUnlock()
	// 返回浅拷贝防止外部修改
	snapshot := make(UploadConfigMap, len(m.configs))
	for k, v := range m.configs {
		cp := *v
		snapshot[k] = &cp
	}
	return snapshot
}

// GetByCode 根据 code 查询当前节点的配置
func GetByCode(code string) (*UploadConfig, bool) {
	m := GetManager()
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := buildKey(m.uploadSource, code)
	v, ok := m.configs[key]
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}

// GetEnabled 获取当前节点所有 enabled=1 的配置列表
func GetEnabled() []*UploadConfig {
	m := GetManager()
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*UploadConfig
	for _, v := range m.configs {
		if v.Enabled == 1 {
			cp := *v
			result = append(result, &cp)
		}
	}
	return result
}

// ─── 变更通知钩子 ────────────────────────────────────────────────────────────

// OnChange 注册配置变更回调，每次热更新检测到变化时触发
// 回调中拿到的是最新配置快照（已加锁保护的副本）
func OnChange(fn func(UploadConfigMap)) {
	m := GetManager()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = append(m.onChange, fn)
}

// ─── 手动刷新 ────────────────────────────────────────────────────────────────

// ForceReload 手动触发一次强制刷新（可用于运维接口）
func ForceReload() error {
	return GetManager().load()
}

// Stop 停止后台热更新（优雅关闭时调用）
func Stop() {
	if globalManager != nil && globalManager.cancel != nil {
		globalManager.cancel()
	}
}

// ─── 内部实现 ────────────────────────────────────────────────────────────────

// watch 后台轮询主循环
func (m *configManager) watch(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	log.Printf("[config] hot-reload watcher started, interval=%s, source=%s", m.interval, m.uploadSource)
	for {
		select {
		case <-ctx.Done():
			log.Println("[config] hot-reload watcher stopped")
			return
		case <-ticker.C:
			if err := m.loadIfChanged(); err != nil {
				log.Printf("[config] reload error: %v", err)
			}
		}
	}
}

// loadIfChanged 检测 updated_at 最大值，有变化才拉取全量
func (m *configManager) loadIfChanged() error {
	var maxUpdatedAt time.Time
	err := m.db.Model(&UploadConfig{}).
		Where("upload_source = ?", m.uploadSource).
		Select("MAX(updated_at)").
		Scan(&maxUpdatedAt).Error
	if err != nil {
		return fmt.Errorf("check updated_at failed: %w", err)
	}

	m.mu.RLock()
	lastCheck := m.lastCheck
	m.mu.RUnlock()

	if !maxUpdatedAt.After(lastCheck) {
		return nil // 无变化，跳过
	}
	return m.load()
}

// load 从 MySQL 全量加载当前节点配置
func (m *configManager) load() error {
	var rows []UploadConfig
	err := m.db.Where("upload_source = ?", m.uploadSource).Find(&rows).Error
	if err != nil {
		return fmt.Errorf("load upload_config failed: %w", err)
	}

	newMap := make(UploadConfigMap, len(rows))
	for i := range rows {
		key := buildKey(rows[i].UploadSource, rows[i].Code)
		newMap[key] = &rows[i]
	}

	m.mu.Lock()
	m.configs = newMap
	m.lastCheck = time.Now()
	callbacks := make([]func(UploadConfigMap), len(m.onChange))
	copy(callbacks, m.onChange)
	m.mu.Unlock()

	log.Printf("[config] loaded %d upload_config items for source=%s", len(rows), m.uploadSource)

	// 在锁外执行回调，避免死锁
	if len(callbacks) > 0 {
		snapshot := All()
		for _, fn := range callbacks {
			fn(snapshot)
		}
	}
	return nil
}

func buildKey(uploadSource, code string) string {
	return uploadSource + "::" + code
}
