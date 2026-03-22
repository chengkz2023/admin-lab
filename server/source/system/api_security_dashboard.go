package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiSecurityDashboard = initOrderApiReliableUpload + 1

type initApiSecurityDashboard struct{}

// 注册网安可视化面板 API 种子，确保新库与老库都能补齐接口定义。
func init() {
	system.RegisterInit(initOrderApiSecurityDashboard, &initApiSecurityDashboard{})
}

func (i *initApiSecurityDashboard) InitializerName() string {
	return "sys_apis_security_dashboard"
}

func (i *initApiSecurityDashboard) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiSecurityDashboard) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiSecurityDashboard) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 使用“先查后建”方式做幂等初始化，避免重复写入。
	entities := []sysModel.SysApi{
		{ApiGroup: "SecurityDashboard", Method: "GET", Path: "/securityDashboard/panel", Description: "获取网安可视化面板数据"},
		{ApiGroup: "SecurityDashboard", Method: "GET", Path: "/securityDashboard/drilldown", Description: "获取网安可视化下钻明细"},
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

func (i *initApiSecurityDashboard) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/securityDashboard/panel", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/securityDashboard/drilldown", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
