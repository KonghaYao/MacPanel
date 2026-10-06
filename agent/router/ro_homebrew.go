package router

import (
	v2 "github.com/1Panel-dev/1Panel/agent/app/api/v2"
	"github.com/gin-gonic/gin"
)

type HomebrewRouter struct{}

func (s *HomebrewRouter) InitRouter(Router *gin.RouterGroup) {
	homebrewRouter := Router.Group("homebrew")
	baseApi := v2.ApiGroupApp.BaseApi
	{
		homebrewRouter.GET("/status", baseApi.GetHomebrewStatus)
		homebrewRouter.POST("/search", baseApi.SearchHomebrew)
		homebrewRouter.GET("/list", baseApi.ListHomebrew)
		homebrewRouter.POST("/install", baseApi.InstallHomebrew)
		homebrewRouter.POST("/uninstall", baseApi.UninstallHomebrew)
		homebrewRouter.POST("/upgrade", baseApi.UpgradeHomebrew)
		homebrewRouter.POST("/update", baseApi.UpdateHomebrew)
		homebrewRouter.GET("/mirror", baseApi.GetHomebrewMirror)
		homebrewRouter.POST("/mirror/update", baseApi.UpdateHomebrewMirror)
		homebrewRouter.POST("/doctor", baseApi.DoctorHomebrew)
	}
}
