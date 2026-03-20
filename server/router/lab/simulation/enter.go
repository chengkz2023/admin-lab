package simulation

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	BaseDataSimulationRouter
}

var (
	baseDataSimulationApi = api.ApiGroupApp.LabApiGroup.SimulationApiGroup.BaseDataSimulationApi
)
