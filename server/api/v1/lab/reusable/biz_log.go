package reusable

import (
	"fmt"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	bizlogUtil "github.com/flipped-aurora/gin-vue-admin/server/utils/bizlog"
	"github.com/gin-gonic/gin"
)

type BizLogApi struct{}

func (b *BizLogApi) List(c *gin.Context) {
	module := c.Query("module")
	entityID := c.Query("entityId")
	if module == "" || entityID == "" {
		response.FailWithMessage("module 和 entityId 不能为空", c)
		return
	}

	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if n, err := parsePositiveInt(v); err == nil {
			page = n
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if n, err := parsePositiveInt(v); err == nil {
			pageSize = n
		}
	}

	result := bizLogService.List(module, entityID, page, pageSize)
	response.OkWithDetailed(result, "获取业务日志成功", c)
}

// WriteDemo writes a test log entry. Only used in admin-lab for demo purposes.
func (b *BizLogApi) WriteDemo(c *gin.Context) {
	var req struct {
		Module   string `json:"module"`
		EntityID string `json:"entityId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Module == "" || req.EntityID == "" {
		response.FailWithMessage("module 和 entityId 不能为空", c)
		return
	}
	bizlogUtil.Record(c.Request.Context(), bizlogUtil.Entry{
		Module:       req.Module,
		EntityID:     req.EntityID,
		Action:       "demo_action",
		OperatorID:   utils.GetUserID(c),
		OperatorName: utils.GetUserName(c),
		Remark:       fmt.Sprintf("这是一条测试日志 [%s/%s]", req.Module, req.EntityID),
	})
	response.OkWithMessage("测试日志已写入", c)
}

func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("invalid")
	}
	return n, nil
}
