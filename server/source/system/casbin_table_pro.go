package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbinTablePro = initOrderCasbinSecurityDashboard + 1

type initCasbinTablePro struct{}

func init() {
	system.RegisterInit(initOrderCasbinTablePro, &initCasbinTablePro{})
}

func (i *initCasbinTablePro) InitializerName() string {
	return "casbin_table_pro"
}

func (i *initCasbinTablePro) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initCasbinTablePro) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initCasbinTablePro) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/tablePro/page", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/tablePro/export", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/tablePro/page", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/tablePro/export", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/tablePro/page", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/tablePro/export", V2: "POST"},
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

func (i *initCasbinTablePro) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/tablePro/page", V2: "POST"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/tablePro/export", V2: "POST"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
