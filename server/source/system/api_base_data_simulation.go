package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiBaseDataSimulation = initOrderApi + 2

type initApiBaseDataSimulation struct{}

func init() {
	system.RegisterInit(initOrderApiBaseDataSimulation, &initApiBaseDataSimulation{})
}

func (i *initApiBaseDataSimulation) InitializerName() string {
	return "sys_apis_base_data_simulation"
}

func (i *initApiBaseDataSimulation) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiBaseDataSimulation) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiBaseDataSimulation) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysApi{
		{ApiGroup: "基础数据导入导出仿真", Method: "GET", Path: "/baseDataSimulation/templates", Description: "获取仿真模板列表"},
		{ApiGroup: "基础数据导入导出仿真", Method: "GET", Path: "/baseDataSimulation/template", Description: "下载仿真模板"},
		{ApiGroup: "基础数据导入导出仿真", Method: "GET", Path: "/baseDataSimulation/export", Description: "导出仿真数据"},
		{ApiGroup: "基础数据导入导出仿真", Method: "POST", Path: "/baseDataSimulation/import", Description: "导入仿真数据（占位）"},
	}

	for _, entity := range entities {
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

func (i *initApiBaseDataSimulation) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/baseDataSimulation/templates", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/baseDataSimulation/import", "POST").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
