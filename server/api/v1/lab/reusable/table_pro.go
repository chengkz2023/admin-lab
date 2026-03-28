package reusable

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TableProApi struct{}

func (t *TableProApi) GetPage(c *gin.Context) {
	var query reusableReq.TableProPageQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.FailWithMessage("invalid query payload", c)
		return
	}
	response.OkWithDetailed(tableProService.Page(query), "get table pro page success", c)
}

func (t *TableProApi) Export(c *gin.Context) {
	var query reusableReq.TableProExportQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.FailWithMessage("invalid export payload", c)
		return
	}

	data, fileName, err := tableProService.Export(query)
	if err != nil {
		global.GVA_LOG.Error("export table pro failed", zap.Error(err))
		response.FailWithMessage("export table pro failed", c)
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
