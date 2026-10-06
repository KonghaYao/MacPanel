package router

import (
	v2 "github.com/1Panel-dev/1Panel/agent/app/api/v2"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/i18n"
	"github.com/1Panel-dev/1Panel/agent/middleware"
	rou "github.com/1Panel-dev/1Panel/agent/router"
	pkgmiddleware "github.com/1Panel-dev/1Panel/pkg/platform/middleware"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func Routers() *gin.Engine {
	Router = gin.Default()
	Router.Use(pkgmiddleware.SlowRequest(pkgmiddleware.SlowRequestConfig{
		Component: "agent",
		Logf:      global.LOG.Warnf,
	}))
	Router.Use(i18n.UseI18n())

	PrivateGroup := Router.Group("/api/v2")
	if !global.IsMaster {
		PrivateGroup.Use(middleware.Certificate())
	}
	PrivateGroup.Use(middleware.OperationResolveMeta())
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}
	PrivateGroup.GET("/health/check", v2.ApiGroupApp.BaseApi.CheckHealth)

	return Router
}
