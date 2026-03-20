package example

import "time"

type ExaReliableUploadLog struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskCode   string    `json:"taskCode" gorm:"column:task_code;type:varchar(64);not null;index:idx_scan,priority:1"`
	TimeStart  time.Time `json:"timeStart" gorm:"column:time_start;not null;index:idx_scan,priority:3"`
	TimeEnd    time.Time `json:"timeEnd" gorm:"column:time_end;not null"`
	FileName   string    `json:"fileName" gorm:"column:file_name;type:varchar(255);not null;uniqueIndex:uk_file_name"`
	BizKey     string    `json:"bizKey" gorm:"column:biz_key;type:varchar(128)"`
	MetaJSON   string    `json:"metaJSON" gorm:"column:meta_json;type:text"`
	Status     uint8     `json:"status" gorm:"column:status;not null;default:0;index:idx_scan,priority:2"`
	BackupPath string    `json:"backupPath" gorm:"column:backup_path;type:varchar(512)"`
	RetryCount int       `json:"retryCount" gorm:"column:retry_count;not null;default:0"`
	ErrMsg     string    `json:"errMsg" gorm:"column:err_msg;type:varchar(1024)"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (ExaReliableUploadLog) TableName() string { return "uploadlog" }

type ExaReliableUploadBigTaskInstance struct {
	ID              int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskCode        string     `json:"taskCode" gorm:"column:task_code;type:varchar(64);not null;uniqueIndex:uk_task_window,priority:1;index:idx_status,priority:1"`
	WindowStart     time.Time  `json:"windowStart" gorm:"column:window_start;not null;uniqueIndex:uk_task_window,priority:2"`
	WindowEnd       time.Time  `json:"windowEnd" gorm:"column:window_end;not null;uniqueIndex:uk_task_window,priority:3"`
	Status          uint8      `json:"status" gorm:"column:status;not null;default:0;index:idx_status,priority:2"`
	TotalBatches    int        `json:"totalBatches" gorm:"column:total_batches"`
	UploadedBatches int        `json:"uploadedBatches" gorm:"column:uploaded_batches;not null;default:0"`
	TotalRecords    int        `json:"totalRecords" gorm:"column:total_records"`
	StartedAt       time.Time  `json:"startedAt" gorm:"column:started_at"`
	FinishedAt      *time.Time `json:"finishedAt" gorm:"column:finished_at"`
}

func (ExaReliableUploadBigTaskInstance) TableName() string { return "big_task_instance" }

type ExaReliableUploadBigTaskBatch struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	InstanceID  int64     `json:"instanceId" gorm:"column:instance_id;not null;uniqueIndex:uk_instance_batch,priority:1;index:idx_instance_status,priority:1"`
	BatchIndex  int       `json:"batchIndex" gorm:"column:batch_index;not null;uniqueIndex:uk_instance_batch,priority:2;index:idx_instance_status,priority:3"`
	FileName    string    `json:"fileName" gorm:"column:file_name;type:varchar(255);not null;uniqueIndex:uk_file_name"`
	RecordCount int       `json:"recordCount" gorm:"column:record_count;not null;default:0"`
	BizKey      string    `json:"bizKey" gorm:"column:biz_key;type:varchar(128)"`
	MetaJSON    string    `json:"metaJSON" gorm:"column:meta_json;type:text"`
	BackupPath  string    `json:"backupPath" gorm:"column:backup_path;type:varchar(512)"`
	Status      uint8     `json:"status" gorm:"column:status;not null;default:0;index:idx_instance_status,priority:2"`
	RetryCount  int       `json:"retryCount" gorm:"column:retry_count;not null;default:0"`
	ErrMsg      string    `json:"errMsg" gorm:"column:err_msg;type:varchar(1024)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (ExaReliableUploadBigTaskBatch) TableName() string { return "big_task_batch" }

type ExaReliableUploadBizTaskInstance struct {
	ID              int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskCode        string     `json:"taskCode" gorm:"column:task_code;type:varchar(64);not null;uniqueIndex:uk_task_trigger,priority:1;index:idx_status,priority:1"`
	TriggerKey      string     `json:"triggerKey" gorm:"column:trigger_key;type:varchar(128);not null;uniqueIndex:uk_task_trigger,priority:2"`
	TriggerPayload  string     `json:"triggerPayload" gorm:"column:trigger_payload;type:text"`
	Status          uint8      `json:"status" gorm:"column:status;not null;default:0;index:idx_status,priority:2"`
	TotalBatches    int        `json:"totalBatches" gorm:"column:total_batches"`
	UploadedBatches int        `json:"uploadedBatches" gorm:"column:uploaded_batches;not null;default:0"`
	TotalRecords    int        `json:"totalRecords" gorm:"column:total_records"`
	StartedAt       time.Time  `json:"startedAt" gorm:"column:started_at"`
	FinishedAt      *time.Time `json:"finishedAt" gorm:"column:finished_at"`
}

func (ExaReliableUploadBizTaskInstance) TableName() string { return "biz_task_instance" }

type ExaReliableUploadBizTaskBatch struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	InstanceID  int64     `json:"instanceId" gorm:"column:instance_id;not null;uniqueIndex:uk_instance_batch,priority:1;index:idx_instance_status,priority:1"`
	BatchIndex  int       `json:"batchIndex" gorm:"column:batch_index;not null;uniqueIndex:uk_instance_batch,priority:2;index:idx_instance_status,priority:3"`
	FileName    string    `json:"fileName" gorm:"column:file_name;type:varchar(255);not null;uniqueIndex:uk_file_name"`
	RecordCount int       `json:"recordCount" gorm:"column:record_count;not null;default:0"`
	BizKey      string    `json:"bizKey" gorm:"column:biz_key;type:varchar(128)"`
	MetaJSON    string    `json:"metaJSON" gorm:"column:meta_json;type:text"`
	BackupPath  string    `json:"backupPath" gorm:"column:backup_path;type:varchar(512)"`
	Status      uint8     `json:"status" gorm:"column:status;not null;default:0;index:idx_instance_status,priority:2"`
	RetryCount  int       `json:"retryCount" gorm:"column:retry_count;not null;default:0"`
	ErrMsg      string    `json:"errMsg" gorm:"column:err_msg;type:varchar(1024)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (ExaReliableUploadBizTaskBatch) TableName() string { return "biz_task_batch" }
