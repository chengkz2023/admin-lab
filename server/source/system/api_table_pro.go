package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiTablePro = initOrderApiSecurityDashboard + 1

type initApiTablePro struct{}

func init() {
	system.RegisterInit(initOrderApiTablePro, &initApiTablePro{})
}

func (i *initApiTablePro) InitializerName() string {
	return "sys_apis_table_pro"
}

func (i *initApiTablePro) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiTablePro) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiTablePro) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysApi{
		{ApiGroup: "TablePro", Method: "POST", Path: "/tablePro/page", Description: "Get table pro page"},
		{ApiGroup: "TablePro", Method: "POST", Path: "/tablePro/export", Description: "Export table pro data"},
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

func (i *initApiTablePro) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/tablePro/page", "POST").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/tablePro/export", "POST").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
