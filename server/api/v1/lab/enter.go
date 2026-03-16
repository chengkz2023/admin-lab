package lab

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/lab/componentdemo"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/lab/reusable"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/lab/simulation"
)

type ApiGroup struct {
	SimulationApiGroup    simulation.ApiGroup
	ComponentDemoApiGroup componentdemo.ApiGroup
	ReusableApiGroup      reusable.ApiGroup
}
