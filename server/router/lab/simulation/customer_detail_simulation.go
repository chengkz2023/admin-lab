package simulation

import "github.com/gin-gonic/gin"

type CustomerDetailSimulationRouter struct{}

func (c *CustomerDetailSimulationRouter) InitCustomerDetailSimulationRouter(Router *gin.RouterGroup) {
	group := Router.Group("customerDetailSimulation")
	{
		group.GET("detail", customerDetailSimulationApi.GetDetail)
	}
}
