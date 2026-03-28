package reusable

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	ExcelIORouter
	ReliableUploadRouter
	SecurityDashboardRouter
	TableProRouter
}

var (
	excelIOApi           = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ExcelIOApi
	reliableUploadApi    = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ReliableUploadApi
	securityDashboardApi = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.SecurityDashboardApi
	tableProApi          = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.TableProApi
)
