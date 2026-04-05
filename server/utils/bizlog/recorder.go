package bizlog

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
	"go.uber.org/zap"
)

// Entry is a business log record.
type Entry struct {
	Module       string
	EntityID     string
	Action       string
	OperatorID   uint
	OperatorName string
	Remark       string
}

// Record writes a business log entry asynchronously.
// On failure, only logs a warn — never blocks or panics.
func Record(ctx context.Context, e Entry) {
	go func() {
		row := reusableModel.BizLog{
			Module:       e.Module,
			EntityID:     e.EntityID,
			Action:       e.Action,
			OperatorID:   e.OperatorID,
			OperatorName: e.OperatorName,
			Remark:       e.Remark,
		}
		if err := global.GVA_DB.WithContext(ctx).Create(&row).Error; err != nil {
			global.GVA_LOG.Warn("bizlog: write failed", zap.String("module", e.Module), zap.String("entityId", e.EntityID), zap.Error(err))
		}
	}()
}
