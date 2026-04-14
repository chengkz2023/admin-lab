package simulation

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CustomerDetailSimulationApi struct{}

func (c *CustomerDetailSimulationApi) GetDetail(ctx *gin.Context) {
	response.OkWithDetailed(customerDetailSimulationService.GetDetail(), "customer detail simulation loaded", ctx)
}
