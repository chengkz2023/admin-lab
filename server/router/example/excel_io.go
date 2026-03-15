package example

import (
	"github.com/gin-gonic/gin"
)

type ExcelIORouter struct{}

func (e *ExcelIORouter) InitExcelIORouter(Router *gin.RouterGroup) {
	excelIORouter := Router.Group("excelIO")
	{
		excelIORouter.GET("templates", excelIOApi.ListTemplates)
		excelIORouter.GET("template", excelIOApi.DownloadTemplate)
		excelIORouter.GET("export", excelIOApi.ExportSample)
		excelIORouter.POST("import", excelIOApi.ImportExcel)
	}
}
