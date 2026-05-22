package config

import "time"

// UploadConfig 对应 cu_upload_config 表
type UploadConfig struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement"`
	Code            string    `gorm:"column:code;not null"`            // x1规范1-67上报文件 x2寻址码
	TaskType        int       `gorm:"column:task_type;not null"`       // 1-分钟实时 2-定时大任务 3-其他指令或特殊上报
	Mode            int       `gorm:"column:mode;default:0"`           // 0-非续传 1-续传
	Interface       int8      `gorm:"column:interface;default:1"`      // 1-x1 2-x2
	IntervalMinutes int       `gorm:"column:interval_minutes;not null"` // 分钟实时间隔，默认1
	DelaySeconds    int       `gorm:"column:delay_seconds;not null"`   // 实时延时，单位秒
	BatchSize       *int      `gorm:"column:batch_size"`               // 分批或分页大小（可为空）
	MaxRetry        int       `gorm:"column:max_retry;default:3"`      // 重试次数
	Enabled         int8      `gorm:"column:enabled;default:1"`        // 0-关 1-开
	Remark          *string   `gorm:"column:remark"`                   // 备注（可为空）
	UploadSource    string    `gorm:"column:upload_source;not null"`   // 上报源，全局唯一标志
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"` // 用于变更检测（建议在表中加此字段）
}

func (UploadConfig) TableName() string {
	return "cu_upload_config"
}
