package router

import (
	v2 "github.com/1Panel-dev/1Panel/agent/app/api/v2"
	"github.com/gin-gonic/gin"
)

type MirrorRouter struct{}

func (s *MirrorRouter) InitRouter(Router *gin.RouterGroup) {
	mirrorRouter := Router.Group("mirrors")
	baseApi := v2.ApiGroupApp.BaseApi
	{
		mirrorRouter.GET("", baseApi.ListMirrors)
		mirrorRouter.POST("/apply", baseApi.ApplyMirror)
	}
}
