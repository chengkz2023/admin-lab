package reusable

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ExcelIOApi
}

var (
	excelIOService = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ExcelIOService
)
