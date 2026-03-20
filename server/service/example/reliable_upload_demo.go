package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	exampleModel "github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/reliableupload"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReliableUploadExampleService struct{}

type ReliableUploadDemoRuntime struct {
	Engine   *reliableupload.Engine
	Reporter *ReliableUploadDemoReporter
}

func (s *ReliableUploadExampleService) AutoMigrateMySQLDemoTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&exampleModel.ExaReliableUploadLog{},
		&exampleModel.ExaReliableUploadBigTaskInstance{},
		&exampleModel.ExaReliableUploadBigTaskBatch{},
		&exampleModel.ExaReliableUploadBizTaskInstance{},
		&exampleModel.ExaReliableUploadBizTaskBatch{},
	)
}

func (s *ReliableUploadExampleService) BuildMySQLDemoRuntime(db *gorm.DB, backupRoot string) (*ReliableUploadDemoRuntime, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	registry := reliableupload.NewRegistry()
	ds := &ReliableUploadDemoDataSource{}
	rp := NewReliableUploadDemoReporter()

	registry.RegisterDataSource("order_minute", ds)
	registry.RegisterReporter("order_minute", rp)
	registry.RegisterDataSource("order_big", ds)
	registry.RegisterReporter("order_big", rp)
	registry.RegisterDataSource("order_biz", ds)
	registry.RegisterReporter("order_biz", rp)

	engine := reliableupload.NewEngine(
		registry,
		newReliableUploadDemoConfigRepo(),
		&mysqlUploadLogRepo{db: db},
		&mysqlBigRepo{db: db},
		&mysqlBizRepo{db: db},
		reliableupload.NewFSBackupStore(backupRoot),
	)

	return &ReliableUploadDemoRuntime{Engine: engine, Reporter: rp}, nil
}

type ReliableUploadDemoDataSource struct{}

func (d *ReliableUploadDemoDataSource) CountChunks(ctx context.Context, cfg reliableupload.TaskConfig, _, _ time.Time) (int, error) {
	if cfg.TaskType == reliableupload.TaskTypeBig {
		return 2, nil
	}
	if cfg.TaskType == reliableupload.TaskTypeBiz {
		if trigger, ok := reliableupload.BizTriggerFromContext(ctx); ok && trigger.Key != "" {
			return 2, nil
		}
		return 1, nil
	}
	return 1, nil
}

func (d *ReliableUploadDemoDataSource) FetchChunk(ctx context.Context, cfg reliableupload.TaskConfig, start, _ time.Time, index int) (reliableupload.Chunk, error) {
	if cfg.TaskType == reliableupload.TaskTypeBig {
		return reliableupload.Chunk{
			Data:        []byte(fmt.Sprintf("%s chunk-%d %s", cfg.TaskCode, index, start.Format("2006-01-02"))),
			RecordCount: 100,
			BizKey:      fmt.Sprintf("biz-%d", index),
			Meta:        map[string]string{"channel": "demo", "bucket": fmt.Sprintf("%d", index)},
		}, nil
	}
	if cfg.TaskType == reliableupload.TaskTypeBiz {
		trigger, _ := reliableupload.BizTriggerFromContext(ctx)
		return reliableupload.Chunk{
			Data:        []byte(fmt.Sprintf("%s trigger=%s payload=%s chunk=%d", cfg.TaskCode, trigger.Key, trigger.Payload, index)),
			RecordCount: 50,
			BizKey:      trigger.Key,
			Meta:        map[string]string{"trigger_key": trigger.Key, "kind": "biz"},
		}, nil
	}
	return reliableupload.Chunk{
		Data: []byte(fmt.Sprintf("%s %s", cfg.TaskCode, start.Format(time.RFC3339))),
		Meta: map[string]string{"kind": "minute"},
	}, nil
}

type ReliableUploadDemoReporter struct {
	mu       sync.Mutex
	uploaded map[string]struct{}
}

func NewReliableUploadDemoReporter() *ReliableUploadDemoReporter {
	return &ReliableUploadDemoReporter{uploaded: map[string]struct{}{}}
}

func (r *ReliableUploadDemoReporter) Upload(_ context.Context, cfg reliableupload.TaskConfig, item reliableupload.UploadItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.uploaded[cfg.TaskCode+"/"+item.FileName] = struct{}{}
	return nil
}

func (r *ReliableUploadDemoReporter) UploadedFiles() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, 0, len(r.uploaded))
	for key := range r.uploaded {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

type reliableUploadDemoConfigRepo struct {
	m map[string]reliableupload.TaskConfig
}

func newReliableUploadDemoConfigRepo() *reliableUploadDemoConfigRepo {
	return &reliableUploadDemoConfigRepo{m: map[string]reliableupload.TaskConfig{
		"order_minute": {TaskCode: "order_minute", TaskType: reliableupload.TaskTypeMinute, IntervalMinutes: 5, DelaySeconds: 60, BatchSize: 500, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order", Enabled: true},
		"order_big":    {TaskCode: "order_big", TaskType: reliableupload.TaskTypeBig, BatchSize: 2000, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order_big", Enabled: true},
		"order_biz":    {TaskCode: "order_biz", TaskType: reliableupload.TaskTypeBiz, BatchSize: 2000, MaxRetry: 3, SFTPSubdir: "/remote/order", FilePrefix: "order_biz", Enabled: true},
	}}
}

func (r *reliableUploadDemoConfigRepo) FindEnabledByType(_ context.Context, typ reliableupload.TaskType) ([]reliableupload.TaskConfig, error) {
	var out []reliableupload.TaskConfig
	for _, cfg := range r.m {
		if cfg.Enabled && cfg.TaskType == typ {
			out = append(out, cfg)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TaskCode < out[j].TaskCode })
	return out, nil
}

func (r *reliableUploadDemoConfigRepo) Get(_ context.Context, taskCode string) (reliableupload.TaskConfig, error) {
	cfg, ok := r.m[taskCode]
	if !ok {
		return reliableupload.TaskConfig{}, fmt.Errorf("task not found: %s", taskCode)
	}
	return cfg, nil
}

type mysqlUploadLogRepo struct{ db *gorm.DB }
type mysqlBigRepo struct{ db *gorm.DB }
type mysqlBizRepo struct{ db *gorm.DB }

func (r *mysqlUploadLogRepo) ExistsByTaskAndTimeRange(ctx context.Context, taskCode string, start, end time.Time) (bool, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("task_code = ? AND time_start = ? AND time_end = ?", taskCode, start, end).Count(&cnt).Error
	return cnt > 0, err
}

func (r *mysqlUploadLogRepo) Create(ctx context.Context, log reliableupload.UploadLog) error {
	return r.db.WithContext(ctx).Create(&exampleModel.ExaReliableUploadLog{TaskCode: log.TaskCode, TimeStart: log.TimeStart, TimeEnd: log.TimeEnd, FileName: log.FileName, BizKey: log.BizKey, MetaJSON: log.MetaJSON, Status: uint8(log.Status), BackupPath: log.BackupPath, RetryCount: log.RetryCount, ErrMsg: log.ErrMsg, CreatedAt: log.CreatedAt, UpdatedAt: log.UpdatedAt}).Error
}

func (r *mysqlUploadLogRepo) FindDistinctPendingTaskCodes(ctx context.Context) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("status = ?", uint8(reliableupload.StatusPending)).Distinct("task_code").Order("task_code ASC").Pluck("task_code", &codes).Error
	return codes, err
}

func (r *mysqlUploadLogRepo) FindPendingByCode(ctx context.Context, taskCode string, maxRetry, limit int) ([]reliableupload.UploadLog, error) {
	var rows []exampleModel.ExaReliableUploadLog
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("task_code = ? AND status = ? AND retry_count <= ?", taskCode, uint8(reliableupload.StatusPending), maxRetry).Order("time_start ASC, id ASC").Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]reliableupload.UploadLog, 0, len(rows))
	for _, row := range rows {
		out = append(out, reliableupload.UploadLog{ID: row.ID, TaskCode: row.TaskCode, TimeStart: row.TimeStart, TimeEnd: row.TimeEnd, FileName: row.FileName, BizKey: row.BizKey, MetaJSON: row.MetaJSON, Status: reliableupload.Status(row.Status), BackupPath: row.BackupPath, RetryCount: row.RetryCount, ErrMsg: row.ErrMsg, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt})
	}
	return out, nil
}

func (r *mysqlUploadLogRepo) MarkUploaded(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("id = ?", id).Updates(map[string]any{"status": uint8(reliableupload.StatusUploaded), "updated_at": time.Now()}).Error
}

func (r *mysqlUploadLogRepo) IncrRetry(ctx context.Context, id int64, errMsg string) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("id = ?", id).Updates(map[string]any{"retry_count": gorm.Expr("retry_count + 1"), "err_msg": errMsg, "updated_at": time.Now()}).Error
}

func (r *mysqlUploadLogRepo) GetLastTimeEndByCode(ctx context.Context, taskCode string) (time.Time, bool, error) {
	var row exampleModel.ExaReliableUploadLog
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadLog{}).Where("task_code = ?", taskCode).Order("time_end DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	return row.TimeEnd, true, nil
}

func (r *mysqlBigRepo) GetOrCreateInstance(ctx context.Context, taskCode string, windowStart, windowEnd time.Time) (reliableupload.BigTaskInstance, error) {
	inst := exampleModel.ExaReliableUploadBigTaskInstance{TaskCode: taskCode, WindowStart: windowStart, WindowEnd: windowEnd, Status: uint8(reliableupload.StatusRunning), StartedAt: time.Now()}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "task_code"}, {Name: "window_start"}, {Name: "window_end"}}, DoNothing: true}).Create(&inst).Error; err != nil {
		return reliableupload.BigTaskInstance{}, err
	}
	var got exampleModel.ExaReliableUploadBigTaskInstance
	if err := r.db.WithContext(ctx).Where("task_code = ? AND window_start = ? AND window_end = ?", taskCode, windowStart, windowEnd).First(&got).Error; err != nil {
		return reliableupload.BigTaskInstance{}, err
	}
	return toBigInstance(got), nil
}

func (r *mysqlBigRepo) UpdateProducedMeta(ctx context.Context, instanceID int64, totalBatches, totalRecords int) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskInstance{}).Where("id = ?", instanceID).Updates(map[string]any{"total_batches": totalBatches, "total_records": totalRecords}).Error
}

func (r *mysqlBigRepo) CreateBatch(ctx context.Context, batch reliableupload.BigTaskBatch) error {
	return r.db.WithContext(ctx).Create(&exampleModel.ExaReliableUploadBigTaskBatch{InstanceID: batch.InstanceID, BatchIndex: batch.BatchIndex, FileName: batch.FileName, RecordCount: batch.RecordCount, BizKey: batch.BizKey, MetaJSON: batch.MetaJSON, BackupPath: batch.BackupPath, Status: uint8(batch.Status), RetryCount: batch.RetryCount, ErrMsg: batch.ErrMsg, CreatedAt: batch.CreatedAt, UpdatedAt: batch.UpdatedAt}).Error
}

func (r *mysqlBigRepo) FindRunningInstances(ctx context.Context) ([]reliableupload.BigTaskInstance, error) {
	var rows []exampleModel.ExaReliableUploadBigTaskInstance
	if err := r.db.WithContext(ctx).Where("status = ?", uint8(reliableupload.StatusRunning)).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]reliableupload.BigTaskInstance, 0, len(rows))
	for _, row := range rows {
		out = append(out, toBigInstance(row))
	}
	return out, nil
}

func (r *mysqlBigRepo) CountBatches(ctx context.Context, instanceID int64) (int, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("instance_id = ?", instanceID).Count(&cnt).Error
	return int(cnt), err
}

func (r *mysqlBigRepo) SumBatchRecords(ctx context.Context, instanceID int64) (int, error) {
	var total sql.NullInt64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("instance_id = ?", instanceID).Select("COALESCE(SUM(record_count), 0)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return int(total.Int64), nil
}

func (r *mysqlBigRepo) FindPendingBatches(ctx context.Context, instanceID int64, maxRetry, limit int) ([]reliableupload.BigTaskBatch, error) {
	var rows []exampleModel.ExaReliableUploadBigTaskBatch
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("instance_id = ? AND status = ? AND retry_count <= ?", instanceID, uint8(reliableupload.StatusPending), maxRetry).Order("batch_index ASC").Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]reliableupload.BigTaskBatch, 0, len(rows))
	for _, row := range rows {
		out = append(out, reliableupload.BigTaskBatch{ID: row.ID, InstanceID: row.InstanceID, BatchIndex: row.BatchIndex, FileName: row.FileName, RecordCount: row.RecordCount, BizKey: row.BizKey, MetaJSON: row.MetaJSON, BackupPath: row.BackupPath, Status: reliableupload.Status(row.Status), RetryCount: row.RetryCount, ErrMsg: row.ErrMsg, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt})
	}
	return out, nil
}

func (r *mysqlBigRepo) MarkBatchUploaded(ctx context.Context, batchID int64) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("id = ?", batchID).Updates(map[string]any{"status": uint8(reliableupload.StatusUploaded), "updated_at": time.Now()}).Error
}

func (r *mysqlBigRepo) IncrBatchRetry(ctx context.Context, batchID int64, errMsg string) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("id = ?", batchID).Updates(map[string]any{"retry_count": gorm.Expr("retry_count + 1"), "err_msg": errMsg, "updated_at": time.Now()}).Error
}

func (r *mysqlBigRepo) CountUploadedBatches(ctx context.Context, instanceID int64) (int, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskBatch{}).Where("instance_id = ? AND status = ?", instanceID, uint8(reliableupload.StatusUploaded)).Count(&cnt).Error
	return int(cnt), err
}

func (r *mysqlBigRepo) MarkInstanceCompleted(ctx context.Context, instanceID int64, finishedAt time.Time) error {
	uploaded, err := r.CountUploadedBatches(ctx, instanceID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBigTaskInstance{}).Where("id = ?", instanceID).Updates(map[string]any{"status": uint8(reliableupload.StatusUploaded), "uploaded_batches": uploaded, "finished_at": finishedAt}).Error
}

func (r *mysqlBizRepo) GetOrCreateInstance(ctx context.Context, taskCode, triggerKey, triggerPayload string) (reliableupload.BizTaskInstance, error) {
	inst := exampleModel.ExaReliableUploadBizTaskInstance{TaskCode: taskCode, TriggerKey: triggerKey, TriggerPayload: triggerPayload, Status: uint8(reliableupload.StatusRunning), StartedAt: time.Now()}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "task_code"}, {Name: "trigger_key"}}, DoNothing: true}).Create(&inst).Error; err != nil {
		return reliableupload.BizTaskInstance{}, err
	}
	var got exampleModel.ExaReliableUploadBizTaskInstance
	if err := r.db.WithContext(ctx).Where("task_code = ? AND trigger_key = ?", taskCode, triggerKey).First(&got).Error; err != nil {
		return reliableupload.BizTaskInstance{}, err
	}
	return toBizInstance(got), nil
}

func (r *mysqlBizRepo) UpdateProducedMeta(ctx context.Context, instanceID int64, totalBatches, totalRecords int) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskInstance{}).Where("id = ?", instanceID).Updates(map[string]any{"total_batches": totalBatches, "total_records": totalRecords}).Error
}

func (r *mysqlBizRepo) CreateBatch(ctx context.Context, batch reliableupload.BizTaskBatch) error {
	return r.db.WithContext(ctx).Create(&exampleModel.ExaReliableUploadBizTaskBatch{InstanceID: batch.InstanceID, BatchIndex: batch.BatchIndex, FileName: batch.FileName, RecordCount: batch.RecordCount, BizKey: batch.BizKey, MetaJSON: batch.MetaJSON, BackupPath: batch.BackupPath, Status: uint8(batch.Status), RetryCount: batch.RetryCount, ErrMsg: batch.ErrMsg, CreatedAt: batch.CreatedAt, UpdatedAt: batch.UpdatedAt}).Error
}

func (r *mysqlBizRepo) FindRunningInstances(ctx context.Context) ([]reliableupload.BizTaskInstance, error) {
	var rows []exampleModel.ExaReliableUploadBizTaskInstance
	if err := r.db.WithContext(ctx).Where("status = ?", uint8(reliableupload.StatusRunning)).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]reliableupload.BizTaskInstance, 0, len(rows))
	for _, row := range rows {
		out = append(out, toBizInstance(row))
	}
	return out, nil
}

func (r *mysqlBizRepo) CountBatches(ctx context.Context, instanceID int64) (int, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("instance_id = ?", instanceID).Count(&cnt).Error
	return int(cnt), err
}

func (r *mysqlBizRepo) SumBatchRecords(ctx context.Context, instanceID int64) (int, error) {
	var total sql.NullInt64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("instance_id = ?", instanceID).Select("COALESCE(SUM(record_count), 0)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return int(total.Int64), nil
}

func (r *mysqlBizRepo) FindPendingBatches(ctx context.Context, instanceID int64, maxRetry, limit int) ([]reliableupload.BizTaskBatch, error) {
	var rows []exampleModel.ExaReliableUploadBizTaskBatch
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("instance_id = ? AND status = ? AND retry_count <= ?", instanceID, uint8(reliableupload.StatusPending), maxRetry).Order("batch_index ASC").Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]reliableupload.BizTaskBatch, 0, len(rows))
	for _, row := range rows {
		out = append(out, reliableupload.BizTaskBatch{ID: row.ID, InstanceID: row.InstanceID, BatchIndex: row.BatchIndex, FileName: row.FileName, RecordCount: row.RecordCount, BizKey: row.BizKey, MetaJSON: row.MetaJSON, BackupPath: row.BackupPath, Status: reliableupload.Status(row.Status), RetryCount: row.RetryCount, ErrMsg: row.ErrMsg, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt})
	}
	return out, nil
}

func (r *mysqlBizRepo) MarkBatchUploaded(ctx context.Context, batchID int64) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("id = ?", batchID).Updates(map[string]any{"status": uint8(reliableupload.StatusUploaded), "updated_at": time.Now()}).Error
}

func (r *mysqlBizRepo) IncrBatchRetry(ctx context.Context, batchID int64, errMsg string) error {
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("id = ?", batchID).Updates(map[string]any{"retry_count": gorm.Expr("retry_count + 1"), "err_msg": errMsg, "updated_at": time.Now()}).Error
}

func (r *mysqlBizRepo) CountUploadedBatches(ctx context.Context, instanceID int64) (int, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskBatch{}).Where("instance_id = ? AND status = ?", instanceID, uint8(reliableupload.StatusUploaded)).Count(&cnt).Error
	return int(cnt), err
}

func (r *mysqlBizRepo) MarkInstanceCompleted(ctx context.Context, instanceID int64, finishedAt time.Time) error {
	uploaded, err := r.CountUploadedBatches(ctx, instanceID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&exampleModel.ExaReliableUploadBizTaskInstance{}).Where("id = ?", instanceID).Updates(map[string]any{"status": uint8(reliableupload.StatusUploaded), "uploaded_batches": uploaded, "finished_at": finishedAt}).Error
}

func toBigInstance(m exampleModel.ExaReliableUploadBigTaskInstance) reliableupload.BigTaskInstance {
	return reliableupload.BigTaskInstance{ID: m.ID, TaskCode: m.TaskCode, WindowStart: m.WindowStart, WindowEnd: m.WindowEnd, Status: reliableupload.Status(m.Status), TotalBatches: m.TotalBatches, UploadedBatches: m.UploadedBatches, TotalRecords: m.TotalRecords, StartedAt: m.StartedAt, FinishedAt: m.FinishedAt}
}

func toBizInstance(m exampleModel.ExaReliableUploadBizTaskInstance) reliableupload.BizTaskInstance {
	return reliableupload.BizTaskInstance{ID: m.ID, TaskCode: m.TaskCode, TriggerKey: m.TriggerKey, TriggerPayload: m.TriggerPayload, Status: reliableupload.Status(m.Status), TotalBatches: m.TotalBatches, UploadedBatches: m.UploadedBatches, TotalRecords: m.TotalRecords, StartedAt: m.StartedAt, FinishedAt: m.FinishedAt}
}
