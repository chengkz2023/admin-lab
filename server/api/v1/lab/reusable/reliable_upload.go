package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type ReliableUploadApi struct{}

func (r *ReliableUploadApi) GetProfile(c *gin.Context) {
	response.OkWithDetailed(reliableUploadService.GetProfile(), "获取可靠上报框架资料成功", c)
}
