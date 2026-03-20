package reusable

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ExcelIOApi
	ReliableUploadApi
}

var (
	excelIOService        = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ExcelIOService
	reliableUploadService = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ReliableUploadService
)
