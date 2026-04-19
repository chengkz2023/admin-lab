package lab

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/source/lab/catalog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderLabAPI = system.InitOrderSystem + 27

type initLabAPI struct{}

func init() {
	system.RegisterInit(initOrderLabAPI, &initLabAPI{})
}

func (i *initLabAPI) InitializerName() string {
	return "sys_apis_lab"
}

func (i *initLabAPI) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initLabAPI) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initLabAPI) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	for _, api := range catalog.BuildAPIs() {
		entity := sysModel.SysApi{
			ApiGroup:    api.APIGroup,
			Method:      api.Method,
			Path:        api.Path,
			Description: api.Description,
		}
		var existing sysModel.SysApi
		err := db.Where("path = ? AND method = ?", entity.Path, entity.Method).First(&existing).Error
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

func (i *initLabAPI) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	for _, api := range catalog.BuildAPIs() {
		if errors.Is(db.Where("path = ? AND method = ?", api.Path, api.Method).First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}
