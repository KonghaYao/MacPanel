package v2

import (
	"errors"
	"net/http"

	"github.com/1Panel-dev/1Panel/agent/app/service"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) GetDesktopWallpaper(c *gin.Context) {
	data, err := desktopService.GetCompressedWallpaper()
	if err != nil {
		if errors.Is(err, service.ErrDesktopWallpaperUnsupported) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "private, max-age=300")
	c.Data(http.StatusOK, "image/jpeg", data)
}
