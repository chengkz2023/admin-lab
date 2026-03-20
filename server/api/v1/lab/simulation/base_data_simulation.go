package simulation

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseDataSimulationApi struct{}

func (b *BaseDataSimulationApi) ListTemplates(c *gin.Context) {
	response.OkWithDetailed(baseDataSimulationService.ListTemplates(), "templates loaded", c)
}

func (b *BaseDataSimulationApi) DownloadTemplate(c *gin.Context) {
	data, fileName, err := baseDataSimulationService.DownloadTemplate(c.Query("templateKey"))
	if err != nil {
		global.GVA_LOG.Error("download template failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (b *BaseDataSimulationApi) ExportData(c *gin.Context) {
	data, fileName, err := baseDataSimulationService.ExportData(c.Query("templateKey"))
	if err != nil {
		global.GVA_LOG.Error("export data failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (b *BaseDataSimulationApi) ImportData(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("read upload file failed", c)
		return
	}

	result, err := baseDataSimulationService.ImportData(c.PostForm("templateKey"), header)
	if err != nil {
		global.GVA_LOG.Error("import validation failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "import parsed and validated (not persisted)", c)
}
