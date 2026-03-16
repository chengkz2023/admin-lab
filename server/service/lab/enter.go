package lab

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/lab/componentdemo"
	"github.com/flipped-aurora/gin-vue-admin/server/service/lab/reusable"
	"github.com/flipped-aurora/gin-vue-admin/server/service/lab/simulation"
)

type ServiceGroup struct {
	SimulationServiceGroup    simulation.ServiceGroup
	ComponentDemoServiceGroup componentdemo.ServiceGroup
	ReusableServiceGroup      reusable.ServiceGroup
}
