package v2

import (
	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) GetPlatformCapabilities(c *gin.Context) {
	helper.SuccessWithData(c, capabilities.Current())
}
