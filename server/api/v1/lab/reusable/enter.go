package reusable

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BizLogApi
	ExcelIOApi
	ReliableUploadApi
	SecurityDashboardApi
	TableProApi
}

var (
	bizLogService            = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.BizLogService
	excelIOService           = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ExcelIOService
	reliableUploadService    = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ReliableUploadService
	securityDashboardService = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.SecurityDashboardService
	tableProService          = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.TableProService
)
