package v2

import (
	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) LoadContainerMonitor(c *gin.Context) {
	var req dto.MonitorContainerSearch
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	data, err := monitorService.LoadContainerMonitorData(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) CleanContainerMonitor(c *gin.Context) {
	var req dto.MonitorContainerClean
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := monitorService.CleanContainerMonitor(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}
