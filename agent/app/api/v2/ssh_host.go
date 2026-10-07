package v2

import (
	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/gin-gonic/gin"
)

// @Tags SSH
// @Summary Search SSH client hosts
// @Accept json
// @Param request body dto.SearchWithPage true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /hosts/ssh/hosts/search [post]
func (b *BaseApi) SearchSSHHosts(c *gin.Context) {
	var req dto.SearchWithPage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, data, err := sshHostService.Search(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, dto.PageResult{
		Total: total,
		Items: data,
	})
}

// @Tags SSH
// @Summary Create SSH client host
// @Accept json
// @Param request body dto.SSHHostOperate true "request"
// @Success 200 {object} dto.SSHHostInfo
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /hosts/ssh/hosts [post]
// @x-panel-log {"bodyKeys":["alias","hostName","user"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"添加 SSH 主机 [alias]","formatEN":"add SSH host [alias]"}
func (b *BaseApi) CreateSSHHost(c *gin.Context) {
	var req dto.SSHHostOperate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	info, err := sshHostService.Create(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, info)
}

// @Tags SSH
// @Summary Delete SSH client hosts
// @Accept json
// @Param request body dto.ForceDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /hosts/ssh/hosts/delete [post]
// @x-panel-log {"bodyKeys":["ids"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"删除 SSH 主机","formatEN":"delete SSH hosts"}
func (b *BaseApi) DeleteSSHHosts(c *gin.Context) {
	var req dto.ForceDelete
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := sshHostService.Delete(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags SSH
// @Summary Test SSH client host connection
// @Accept json
// @Param request body dto.SSHHostTest true "request"
// @Success 200 {object} dto.SSHHostTestResult
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /hosts/ssh/hosts/test [post]
func (b *BaseApi) TestSSHHost(c *gin.Context) {
	var req dto.SSHHostTest
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	helper.SuccessWithData(c, dto.SSHHostTestResult{
		Status: sshHostService.Test(req),
	})
}
