package simulation

import "github.com/gin-gonic/gin"

type BaseDataSimulationRouter struct{}

func (b *BaseDataSimulationRouter) InitBaseDataSimulationRouter(Router *gin.RouterGroup) {
	group := Router.Group("baseDataSimulation")
	{
		group.GET("templates", baseDataSimulationApi.ListTemplates)
		group.GET("template", baseDataSimulationApi.DownloadTemplate)
		group.GET("export", baseDataSimulationApi.ExportData)
		group.POST("import", baseDataSimulationApi.ImportData)
	}
}
