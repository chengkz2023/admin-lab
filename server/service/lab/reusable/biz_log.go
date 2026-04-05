package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
)

type BizLogService struct{}

type BizLogListResult struct {
	List  []reusableModel.BizLog `json:"list"`
	Total int64                  `json:"total"`
}

// List 按 module + entityID 倒序分页查询。
func (s *BizLogService) List(module, entityID string, page, pageSize int) BizLogListResult {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var rows []reusableModel.BizLog
	var total int64

	db := global.GVA_DB.Model(&reusableModel.BizLog{}).
		Where("module = ? AND entity_id = ?", module, entityID)

	db.Count(&total)
	db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows)

	return BizLogListResult{List: rows, Total: total}
}
