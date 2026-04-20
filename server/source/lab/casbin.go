package lab

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/source/lab/catalog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderLabCasbin = system.InitOrderSystem + 28

type initLabCasbin struct{}

func init() {
	system.RegisterInit(initOrderLabCasbin, &initLabCasbin{})
}

func (i *initLabCasbin) InitializerName() string {
	return "casbin_lab"
}

func (i *initLabCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initLabCasbin) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initLabCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	for _, policy := range catalog.BuildCasbinPolicies() {
		entity := adapter.CasbinRule{
			Ptype: "p",
			V0:    policy.RoleID,
			V1:    policy.Path,
			V2:    policy.Method,
		}
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

func (i *initLabCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	for _, policy := range catalog.BuildCasbinPolicies() {
		if errors.Is(
			db.Where(adapter.CasbinRule{Ptype: "p", V0: policy.RoleID, V1: policy.Path, V2: policy.Method}).
				First(&adapter.CasbinRule{}).Error,
			gorm.ErrRecordNotFound,
		) {
			return false
		}
	}
	return true
}
