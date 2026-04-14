package simulation

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BaseDataSimulationApi
	CustomerDetailSimulationApi
}

var (
	baseDataSimulationService       = service.ServiceGroupApp.LabServiceGroup.SimulationServiceGroup.BaseDataSimulationService
	customerDetailSimulationService = service.ServiceGroupApp.LabServiceGroup.SimulationServiceGroup.CustomerDetailSimulationService
)
