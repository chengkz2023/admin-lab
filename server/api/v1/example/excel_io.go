package example

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExcelIOApi struct{}

func (e *ExcelIOApi) ListTemplates(c *gin.Context) {
	response.OkWithDetailed(excelIOService.ListTemplates(), "获取 Excel 模板列表成功", c)
}

func (e *ExcelIOApi) DownloadTemplate(c *gin.Context) {
	data, fileName, err := excelIOService.ExportTemplate(c.Query("templateKey"))
	if err != nil {
		global.GVA_LOG.Error("导出 Excel 模板失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", downloadDisposition(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (e *ExcelIOApi) ExportSample(c *gin.Context) {
	data, fileName, err := excelIOService.ExportSample()
	if err != nil {
		global.GVA_LOG.Error("导出 Excel 示例失败", zap.Error(err))
		response.FailWithMessage("导出 Excel 示例失败", c)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", downloadDisposition(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (e *ExcelIOApi) ImportExcel(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("读取上传文件失败", c)
		return
	}

	result, err := excelIOService.ImportExcel(c.PostForm("templateKey"), header)
	if err != nil {
		global.GVA_LOG.Error("导入 Excel 失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "导入解析成功", c)
}

func downloadDisposition(fileName string) string {
	return "attachment; filename=" + fileName
}
