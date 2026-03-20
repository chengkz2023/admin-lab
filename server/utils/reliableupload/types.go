package reliableupload

import "time"

type TaskType uint8

const (
	TaskTypeMinute TaskType = 1
	TaskTypeBig    TaskType = 2
	TaskTypeBiz    TaskType = 3
)

type Status uint8

const (
	StatusPending  Status = 0
	StatusUploaded Status = 1
	StatusFailed   Status = 2
	StatusRunning  Status = 3
)

type TaskConfig struct {
	TaskCode        string
	TaskType        TaskType
	IntervalMinutes int
	DelaySeconds    int
	BatchSize       int
	MaxRetry        int
	SFTPSubdir      string
	FilePrefix      string
	Enabled         bool
}

type Chunk struct {
	Data        []byte
	RecordCount int
	BizKey      string
	Meta        map[string]string
}

type NameContext struct {
	BizKey string
	Meta   map[string]string
}

type UploadItem struct {
	FileName   string
	Data       []byte
	BizKey     string
	Meta       map[string]string
	BackupPath string
}

type UploadLog struct {
	ID         int64
	TaskCode   string
	TimeStart  time.Time
	TimeEnd    time.Time
	FileName   string
	BizKey     string
	MetaJSON   string
	Status     Status
	BackupPath string
	RetryCount int
	ErrMsg     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type BigTaskInstance struct {
	ID              int64
	TaskCode        string
	WindowStart     time.Time
	WindowEnd       time.Time
	Status          Status
	TotalBatches    int
	UploadedBatches int
	TotalRecords    int
	StartedAt       time.Time
	FinishedAt      *time.Time
}

type BigTaskBatch struct {
	ID          int64
	InstanceID  int64
	BatchIndex  int
	FileName    string
	RecordCount int
	BizKey      string
	MetaJSON    string
	BackupPath  string
	Status      Status
	RetryCount  int
	ErrMsg      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type BizTrigger struct {
	Key     string
	Payload string
}

type BizTaskInstance struct {
	ID              int64
	TaskCode        string
	TriggerKey      string
	TriggerPayload  string
	Status          Status
	TotalBatches    int
	UploadedBatches int
	TotalRecords    int
	StartedAt       time.Time
	FinishedAt      *time.Time
}

type BizTaskBatch struct {
	ID          int64
	InstanceID  int64
	BatchIndex  int
	FileName    string
	RecordCount int
	BizKey      string
	MetaJSON    string
	BackupPath  string
	Status      Status
	RetryCount  int
	ErrMsg      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
