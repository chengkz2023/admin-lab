package reusable

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	ExcelIORouter
	ReliableUploadRouter
	SecurityDashboardRouter
}

var (
	excelIOApi           = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ExcelIOApi
	reliableUploadApi    = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ReliableUploadApi
	securityDashboardApi = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.SecurityDashboardApi
)
