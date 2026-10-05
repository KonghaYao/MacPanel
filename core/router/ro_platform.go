package router

import (
	v2 "github.com/1Panel-dev/1Panel/core/app/api/v2"
	"github.com/gin-gonic/gin"
)

type PlatformRouter struct{}

func (s *PlatformRouter) InitRouter(Router *gin.RouterGroup) {
	baseApi := v2.ApiGroupApp.BaseApi
	Router.GET("/capabilities", baseApi.GetPlatformCapabilities)
}
