package simulation

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BaseDataSimulationApi
}

var (
	baseDataSimulationService = service.ServiceGroupApp.LabServiceGroup.SimulationServiceGroup.BaseDataSimulationService
)
