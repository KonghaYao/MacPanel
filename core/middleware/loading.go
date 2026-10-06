package middleware

import (
	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/core/app/repo"
	"github.com/gin-gonic/gin"
)

func GlobalLoading() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/mcp" {
			c.Next()
			return
		}
		settingRepo := repo.NewISettingRepo()
		status, err := settingRepo.GetValueByKey("SystemStatus")
		if err != nil {
			helper.InternalServer(c, err)
			return
		}
		if status != "Free" {
			helper.ErrorWithDetail(c, 407, status, err)
			return
		}
		c.Next()
	}
}
