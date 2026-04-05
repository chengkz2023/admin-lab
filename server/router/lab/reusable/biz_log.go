package reusable

import "github.com/gin-gonic/gin"

type BizLogRouter struct{}

func (r *BizLogRouter) InitBizLogRouter(Router *gin.RouterGroup) {
	bizLogRouter := Router.Group("bizLog")
	{
		bizLogRouter.GET("list", bizLogApi.List)
		bizLogRouter.POST("writeDemo", bizLogApi.WriteDemo)
	}
}
