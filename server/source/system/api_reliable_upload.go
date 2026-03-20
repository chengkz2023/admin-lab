package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiReliableUpload = initOrderApiExcelIO + 1

type initApiReliableUpload struct{}

func init() {
	system.RegisterInit(initOrderApiReliableUpload, &initApiReliableUpload{})
}

func (i *initApiReliableUpload) InitializerName() string {
	return "sys_apis_reliable_upload"
}

func (i *initApiReliableUpload) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiReliableUpload) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiReliableUpload) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entity := sysModel.SysApi{ApiGroup: "ReliableUpload", Method: "GET", Path: "/reliableUpload/profile", Description: "获取可靠上报框架资料"}
	var existing sysModel.SysApi
	err := db.Where("path = ? AND method = ?", entity.Path, entity.Method).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx, db.Create(&entity).Error
	}
	return ctx, err
}

func (i *initApiReliableUpload) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return !errors.Is(db.Where("path = ? AND method = ?", "/reliableUpload/profile", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound)
}
