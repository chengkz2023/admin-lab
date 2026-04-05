package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiBizLog = initOrderApiReliableUpload + 1

type initApiBizLog struct{}

func init() {
	system.RegisterInit(initOrderApiBizLog, &initApiBizLog{})
}

func (i *initApiBizLog) InitializerName() string { return "sys_apis_biz_log" }

func (i *initApiBizLog) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiBizLog) TableCreated(ctx context.Context) bool { return false }

func (i *initApiBizLog) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "BizLog", Method: "GET", Path: "/bizLog/list", Description: "查询业务操作日志"},
		{ApiGroup: "BizLog", Method: "POST", Path: "/bizLog/writeDemo", Description: "写入测试业务日志（仅 admin-lab 演示用）"},
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

func (i *initApiBizLog) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return !errors.Is(db.Where("path = ? AND method = ?", "/bizLog/list", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound)
}
