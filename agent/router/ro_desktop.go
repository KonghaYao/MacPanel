package router

import (
	v2 "github.com/1Panel-dev/1Panel/agent/app/api/v2"
	"github.com/gin-gonic/gin"
)

type DesktopRouter struct{}

func (s *DesktopRouter) InitRouter(Router *gin.RouterGroup) {
	desktopRouter := Router.Group("desktop")
	baseApi := v2.ApiGroupApp.BaseApi
	{
		desktopRouter.GET("/wallpaper", baseApi.GetDesktopWallpaper)
	}
}
