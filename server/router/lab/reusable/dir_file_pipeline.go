package reusable

import "github.com/gin-gonic/gin"

type DirFilePipelineRouter struct{}

func (d *DirFilePipelineRouter) InitDirFilePipelineRouter(Router *gin.RouterGroup) {
	dirFilePipelineRouter := Router.Group("dirPipeline")
	{
		dirFilePipelineRouter.GET("profile", dirFilePipelineApi.GetProfile)
		dirFilePipelineRouter.POST("runOnce", dirFilePipelineApi.RunOnce)
	}
}
