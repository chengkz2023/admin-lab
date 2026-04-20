package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiCustomerDetailSimulation = initOrderApiIgnore + 1

type initApiCustomerDetailSimulation struct{}

func init() {
	system.RegisterInit(initOrderApiCustomerDetailSimulation, &initApiCustomerDetailSimulation{})
}

func (i *initApiCustomerDetailSimulation) InitializerName() string {
	return "sys_apis_customer_detail_simulation"
}

func (i *initApiCustomerDetailSimulation) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiCustomerDetailSimulation) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initApiCustomerDetailSimulation) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysApi{
		{ApiGroup: "客户详情仿真", Method: "GET", Path: "/customerDetailSimulation/detail", Description: "获取客户详情仿真数据"},
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

func (i *initApiCustomerDetailSimulation) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/customerDetailSimulation/detail", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
