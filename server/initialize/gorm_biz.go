package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
)

func bizModel() error {
	db := global.GVA_DB
	return db.AutoMigrate(&reusableModel.BizLog{})
}
