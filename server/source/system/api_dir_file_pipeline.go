package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiDirFilePipeline = initOrderApiIgnore + 1

type initApiDirFilePipeline struct{}

func init() {
	system.RegisterInit(initOrderApiDirFilePipeline, &initApiDirFilePipeline{})
}

func (i *initApiDirFilePipeline) InitializerName() string {
	return "sys_apis_dir_file_pipeline"
}

func (i *initApiDirFilePipeline) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiDirFilePipeline) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiDirFilePipeline) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysApi{
		{ApiGroup: "DirFilePipeline", Method: "GET", Path: "/dirPipeline/profile", Description: "获取目录文件处理流水线介绍"},
		{ApiGroup: "DirFilePipeline", Method: "POST", Path: "/dirPipeline/runOnce", Description: "执行一次目录文件处理流水线"},
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

func (i *initApiDirFilePipeline) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/dirPipeline/profile", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/dirPipeline/runOnce", "POST").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
