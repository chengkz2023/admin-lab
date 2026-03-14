package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"go.uber.org/zap"
)

// EnsureSystemSeedData 在服务启动时补齐系统默认数据。
func EnsureSystemSeedData() {
	if global.GVA_DB == nil {
		return
	}

	if err := service.ServiceGroupApp.SystemServiceGroup.InitDBService.EnsureSystemInitDataOnBoot(global.GVA_DB); err != nil {
		global.GVA_LOG.Error("ensure system seed data failed", zap.Error(err))
	}
}
