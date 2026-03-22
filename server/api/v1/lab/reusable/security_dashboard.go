package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	"github.com/gin-gonic/gin"
)

type SecurityDashboardApi struct{}

// GetPanel 返回可视化面板统一契约数据。
// 该接口面向“复用组件”示例，便于迁移到内网项目直接复用。
func (s *SecurityDashboardApi) GetPanel(c *gin.Context) {
	var query reusableReq.SecurityDashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("查询参数格式错误", c)
		return
	}
	response.OkWithDetailed(securityDashboardService.GetPanel(query), "获取可视化面板数据成功", c)
}

// GetDrilldown 返回图表点击下钻后的事件明细。
// 入参包含图元上下文及当前筛选条件，用于保证上下文一致性。
func (s *SecurityDashboardApi) GetDrilldown(c *gin.Context) {
	var query reusableReq.SecurityDrilldownQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("查询参数格式错误", c)
		return
	}
	response.OkWithDetailed(securityDashboardService.GetDrilldown(query), "获取下钻明细数据成功", c)
}
