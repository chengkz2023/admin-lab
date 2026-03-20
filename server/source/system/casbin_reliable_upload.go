package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbinReliableUpload = initOrderCasbinExcelIO + 1

type initCasbinReliableUpload struct{}

func init() {
	system.RegisterInit(initOrderCasbinReliableUpload, &initCasbinReliableUpload{})
}

func (i *initCasbinReliableUpload) InitializerName() string {
	return "casbin_reliable_upload"
}

func (i *initCasbinReliableUpload) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initCasbinReliableUpload) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initCasbinReliableUpload) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/reliableUpload/profile", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/reliableUpload/profile", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/reliableUpload/profile", V2: "GET"},
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

func (i *initCasbinReliableUpload) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return !errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/reliableUpload/profile", V2: "GET"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound)
}
