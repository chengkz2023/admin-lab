package reusable

import "github.com/gin-gonic/gin"

type SecurityDashboardRouter struct{}

// InitSecurityDashboardRouter 注册网安可视化面板相关接口。
// 统一挂载在 /securityDashboard 下，便于前端集中调用。
func (s *SecurityDashboardRouter) InitSecurityDashboardRouter(Router *gin.RouterGroup) {
	securityDashboardRouter := Router.Group("securityDashboard")
	{
		securityDashboardRouter.GET("panel", securityDashboardApi.GetPanel)
		securityDashboardRouter.GET("drilldown", securityDashboardApi.GetDrilldown)
	}
}
