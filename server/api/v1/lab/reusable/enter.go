package reusable

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ExcelIOApi
	ReliableUploadApi
	SecurityDashboardApi
}

var (
	excelIOService           = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ExcelIOService
	reliableUploadService    = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ReliableUploadService
	securityDashboardService = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.SecurityDashboardService
)
