package reusable

import "time"

// BizLog 业务操作日志，按 module + entity_id 查询。
// 迁入内网时：执行 AutoMigrate 或手动建表；recorder 替换 DB 实例即可。
type BizLog struct {
	ID           uint      `json:"id"           gorm:"primaryKey;autoIncrement"`
	Module       string    `json:"module"       gorm:"size:64;not null;index:idx_module_entity"`
	EntityID     string    `json:"entityId"     gorm:"column:entity_id;size:64;not null;index:idx_module_entity"`
	Action       string    `json:"action"       gorm:"size:64;not null"`
	OperatorID   uint      `json:"operatorId"   gorm:"not null"`
	OperatorName string    `json:"operatorName" gorm:"size:64;not null"`
	Remark       string    `json:"remark"       gorm:"size:512"`
	CreatedAt    time.Time `json:"createdAt"    gorm:"index:idx_created_at"`
}

func (BizLog) TableName() string { return "biz_log" }
