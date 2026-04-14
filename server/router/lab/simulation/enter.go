package simulation

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	BaseDataSimulationRouter
	CustomerDetailSimulationRouter
}

var (
	baseDataSimulationApi       = api.ApiGroupApp.LabApiGroup.SimulationApiGroup.BaseDataSimulationApi
	customerDetailSimulationApi = api.ApiGroupApp.LabApiGroup.SimulationApiGroup.CustomerDetailSimulationApi
)
