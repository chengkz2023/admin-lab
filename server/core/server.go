package core

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

func RunServer() {
	if global.GVA_CONFIG.Redis.Enable {
		initialize.Redis()
	}

	if global.GVA_DB != nil {
		system.LoadAll()
	}

	router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)

	fmt.Printf(`
BoyKing Admin 已启动
版本: %s
访问地址: http://127.0.0.1%s
`, global.Version, address)

	initServer(address, router, 10*time.Minute, 10*time.Minute)
}
