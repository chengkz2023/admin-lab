package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbinBaseDataSimulation = initOrderCasbin + 2

type initCasbinBaseDataSimulation struct{}

func init() {
	system.RegisterInit(initOrderCasbinBaseDataSimulation, &initCasbinBaseDataSimulation{})
}

func (i *initCasbinBaseDataSimulation) InitializerName() string {
	return "casbin_base_data_simulation"
}

func (i *initCasbinBaseDataSimulation) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initCasbinBaseDataSimulation) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initCasbinBaseDataSimulation) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/baseDataSimulation/templates", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/baseDataSimulation/template", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/baseDataSimulation/export", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/baseDataSimulation/import", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/baseDataSimulation/templates", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/baseDataSimulation/template", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/baseDataSimulation/export", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/baseDataSimulation/import", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/templates", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/template", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/export", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/import", V2: "POST"},
	}

	for _, entity := range entities {
		var existing adapter.CasbinRule
		err := db.Where(adapter.CasbinRule{Ptype: entity.Ptype, V0: entity.V0, V1: entity.V1, V2: entity.V2}).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := db.Create(&entity).Error; createErr != nil {
				return ctx, createErr
			}
			continue
		}
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (i *initCasbinBaseDataSimulation) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/templates", V2: "GET"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/baseDataSimulation/import", V2: "POST"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
