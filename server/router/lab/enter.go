package lab

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/lab/componentdemo"
	"github.com/flipped-aurora/gin-vue-admin/server/router/lab/reusable"
	"github.com/flipped-aurora/gin-vue-admin/server/router/lab/simulation"
)

type RouterGroup struct {
	SimulationRouterGroup    simulation.RouterGroup
	ComponentDemoRouterGroup componentdemo.RouterGroup
	ReusableRouterGroup      reusable.RouterGroup
}
