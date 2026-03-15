package initialize

import (
	"net/http"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrPermission
	}
	return f, nil
}

func Routers() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinRecovery(true))
	if gin.Mode() == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example

	engine.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)})

	publicGroup := engine.Group(global.GVA_CONFIG.System.RouterPrefix)
	privateGroup := engine.Group(global.GVA_CONFIG.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	publicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	systemRouter.InitBaseRouter(publicGroup)
	systemRouter.InitInitRouter(publicGroup)
	systemRouter.InitApiRouter(privateGroup, publicGroup)
	systemRouter.InitJwtRouter(privateGroup)
	systemRouter.InitUserRouter(privateGroup)
	systemRouter.InitMenuRouter(privateGroup)
	systemRouter.InitCasbinRouter(privateGroup)
	systemRouter.InitAuthorityRouter(privateGroup)
	systemRouter.InitSysDictionaryRouter(privateGroup)
	systemRouter.InitSysOperationRecordRouter(privateGroup)
	systemRouter.InitSysDictionaryDetailRouter(privateGroup)
	systemRouter.InitAuthorityBtnRouterRouter(privateGroup)
	systemRouter.InitSysParamsRouter(privateGroup, publicGroup)
	exampleRouter.InitFileUploadAndDownloadRouter(privateGroup)
	exampleRouter.InitAttachmentCategoryRouterRouter(privateGroup)
	exampleRouter.InitExcelIORouter(privateGroup)
	initBizRouter(privateGroup, publicGroup)

	global.GVA_ROUTERS = engine.Routes()
	global.GVA_LOG.Info("router register success")
	return engine
}
