package reliableupload

import (
	"context"
	"time"
)

type DataSource interface {
	CountChunks(ctx context.Context, cfg TaskConfig, start, end time.Time) (int, error)
	FetchChunk(ctx context.Context, cfg TaskConfig, start, end time.Time, index int) (Chunk, error)
}

type Reporter interface {
	Upload(ctx context.Context, cfg TaskConfig, item UploadItem) error
}

type BackupStore interface {
	Save(ctx context.Context, taskCode, fileName string, data []byte) (backupPath string, err error)
	Read(ctx context.Context, backupPath string) ([]byte, error)
}

type TaskConfigRepo interface {
	FindEnabledByType(ctx context.Context, typ TaskType) ([]TaskConfig, error)
	Get(ctx context.Context, taskCode string) (TaskConfig, error)
}

type UploadLogRepo interface {
	ExistsByTaskAndTimeRange(ctx context.Context, taskCode string, start, end time.Time) (bool, error)
	Create(ctx context.Context, log UploadLog) error
	FindDistinctPendingTaskCodes(ctx context.Context) ([]string, error)
	FindPendingByCode(ctx context.Context, taskCode string, maxRetry, limit int) ([]UploadLog, error)
	MarkUploaded(ctx context.Context, id int64) error
	IncrRetry(ctx context.Context, id int64, errMsg string) error
	GetLastTimeEndByCode(ctx context.Context, taskCode string) (time.Time, bool, error)
}

type BigTaskRepo interface {
	GetOrCreateInstance(ctx context.Context, taskCode string, windowStart, windowEnd time.Time) (BigTaskInstance, error)
	UpdateProducedMeta(ctx context.Context, instanceID int64, totalBatches, totalRecords int) error
	CreateBatch(ctx context.Context, batch BigTaskBatch) error
	FindRunningInstances(ctx context.Context) ([]BigTaskInstance, error)
	CountBatches(ctx context.Context, instanceID int64) (int, error)
	SumBatchRecords(ctx context.Context, instanceID int64) (int, error)
	FindPendingBatches(ctx context.Context, instanceID int64, maxRetry, limit int) ([]BigTaskBatch, error)
	MarkBatchUploaded(ctx context.Context, batchID int64) error
	IncrBatchRetry(ctx context.Context, batchID int64, errMsg string) error
	CountUploadedBatches(ctx context.Context, instanceID int64) (int, error)
	MarkInstanceCompleted(ctx context.Context, instanceID int64, finishedAt time.Time) error
}

type BizTaskRepo interface {
	GetOrCreateInstance(ctx context.Context, taskCode, triggerKey, triggerPayload string) (BizTaskInstance, error)
	UpdateProducedMeta(ctx context.Context, instanceID int64, totalBatches, totalRecords int) error
	CreateBatch(ctx context.Context, batch BizTaskBatch) error
	FindRunningInstances(ctx context.Context) ([]BizTaskInstance, error)
	CountBatches(ctx context.Context, instanceID int64) (int, error)
	SumBatchRecords(ctx context.Context, instanceID int64) (int, error)
	FindPendingBatches(ctx context.Context, instanceID int64, maxRetry, limit int) ([]BizTaskBatch, error)
	MarkBatchUploaded(ctx context.Context, batchID int64) error
	IncrBatchRetry(ctx context.Context, batchID int64, errMsg string) error
	CountUploadedBatches(ctx context.Context, instanceID int64) (int, error)
	MarkInstanceCompleted(ctx context.Context, instanceID int64, finishedAt time.Time) error
}

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type LoggerFuncs struct {
	InfofFunc  func(format string, args ...any)
	ErrorfFunc func(format string, args ...any)
}

func (l LoggerFuncs) Infof(format string, args ...any) {
	if l.InfofFunc != nil {
		l.InfofFunc(format, args...)
	}
}

func (l LoggerFuncs) Errorf(format string, args ...any) {
	if l.ErrorfFunc != nil {
		l.ErrorfFunc(format, args...)
	}
}

type Clock interface {
	Now() time.Time
}

type FileNamer interface {
	FileName(cfg TaskConfig, windowStart, windowEnd time.Time, batchIndex int, ctx NameContext) string
}

type Reconciler interface {
	Reconcile(ctx context.Context) error
}
