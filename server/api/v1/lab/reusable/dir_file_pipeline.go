package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DirFilePipelineApi struct{}

func (d *DirFilePipelineApi) GetProfile(c *gin.Context) {
	response.OkWithDetailed(dirFilePipelineService.GetProfile(), "目录文件处理流水线介绍加载成功", c)
}

func (d *DirFilePipelineApi) RunOnce(c *gin.Context) {
	var req reusableReq.DirFilePipelineRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数格式不正确", c)
		return
	}

	result, err := dirFilePipelineService.RunOnce(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("目录文件处理流水线执行失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "目录文件处理流水线执行完成", c)
}
