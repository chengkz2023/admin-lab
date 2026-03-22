package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbinSecurityDashboard = initOrderCasbinReliableUpload + 1

type initCasbinSecurityDashboard struct{}

// 注册网安可视化面板 Casbin 种子，保证默认角色可访问新接口。
func init() {
	system.RegisterInit(initOrderCasbinSecurityDashboard, &initCasbinSecurityDashboard{})
}

func (i *initCasbinSecurityDashboard) InitializerName() string {
	return "casbin_security_dashboard"
}

func (i *initCasbinSecurityDashboard) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initCasbinSecurityDashboard) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initCasbinSecurityDashboard) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 覆盖 888 / 8881 / 9528 三个内置角色，采用幂等写入。
	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/securityDashboard/panel", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/securityDashboard/drilldown", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/securityDashboard/panel", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/securityDashboard/drilldown", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/securityDashboard/panel", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/securityDashboard/drilldown", V2: "GET"},
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

func (i *initCasbinSecurityDashboard) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/securityDashboard/panel", V2: "GET"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/securityDashboard/drilldown", V2: "GET"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
