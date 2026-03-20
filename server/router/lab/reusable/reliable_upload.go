package reusable

import "github.com/gin-gonic/gin"

type ReliableUploadRouter struct{}

func (r *ReliableUploadRouter) InitReliableUploadRouter(Router *gin.RouterGroup) {
	reliableUploadRouter := Router.Group("reliableUpload")
	{
		reliableUploadRouter.GET("profile", reliableUploadApi.GetProfile)
	}
}
