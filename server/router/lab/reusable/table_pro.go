package reusable

import "github.com/gin-gonic/gin"

type TableProRouter struct{}

func (t *TableProRouter) InitTableProRouter(Router *gin.RouterGroup) {
	tableProRouter := Router.Group("tablePro")
	{
		tableProRouter.POST("page", tableProApi.GetPage)
		tableProRouter.POST("export", tableProApi.Export)
	}
}
