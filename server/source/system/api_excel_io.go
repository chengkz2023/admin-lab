package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiExcelIO = initOrderApi + 1

type initApiExcelIO struct{}

func init() {
	system.RegisterInit(initOrderApiExcelIO, &initApiExcelIO{})
}

func (i *initApiExcelIO) InitializerName() string {
	return "sys_apis_excel_io"
}

func (i *initApiExcelIO) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiExcelIO) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiExcelIO) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysApi{
		{ApiGroup: "Excel导入导出", Method: "GET", Path: "/excelIO/templates", Description: "获取 Excel 模板列表"},
		{ApiGroup: "Excel导入导出", Method: "GET", Path: "/excelIO/template", Description: "下载 Excel 导入模板"},
		{ApiGroup: "Excel导入导出", Method: "GET", Path: "/excelIO/export", Description: "导出 Excel 示例数据"},
		{ApiGroup: "Excel导入导出", Method: "POST", Path: "/excelIO/import", Description: "导入并解析 Excel"},
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

func (i *initApiExcelIO) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/excelIO/templates", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/excelIO/import", "POST").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
