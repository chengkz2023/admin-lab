package reusable

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	ExcelIORouter
}

var (
	excelIOApi = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ExcelIOApi
)
