package v2

import (
	"net/http"

	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) GetHomebrewStatus(c *gin.Context) {
	status, err := homebrewService.GetStatus()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, status)
}

func (b *BaseApi) SearchHomebrew(c *gin.Context) {
	var req request.HomebrewSearchReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	results, err := homebrewService.Search(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, results)
}

func (b *BaseApi) ListHomebrew(c *gin.Context) {
	var req request.HomebrewListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorWithDetail(c, http.StatusBadRequest, "ErrInvalidParams", err)
		return
	}
	if err := global.VALID.Struct(req); err != nil {
		helper.ErrorWithDetail(c, http.StatusBadRequest, "ErrInvalidParams", err)
		return
	}
	total, records, err := homebrewService.List(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, dto.PageResult{
		Total: total,
		Items: records,
	})
}

func (b *BaseApi) InstallHomebrew(c *gin.Context) {
	var req request.HomebrewPackageReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := homebrewService.Install(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) UninstallHomebrew(c *gin.Context) {
	var req request.HomebrewPackageReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := homebrewService.Uninstall(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) UpgradeHomebrew(c *gin.Context) {
	var req request.HomebrewPackageReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := homebrewService.Upgrade(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) UpdateHomebrew(c *gin.Context) {
	var req request.HomebrewTaskReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := homebrewService.Update(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) GetHomebrewMirror(c *gin.Context) {
	config, err := homebrewService.GetMirror()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, config)
}

func (b *BaseApi) UpdateHomebrewMirror(c *gin.Context) {
	var req request.HomebrewMirrorUpdateReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := homebrewService.UpdateMirror(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) DoctorHomebrew(c *gin.Context) {
	result, err := homebrewService.Doctor()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, result)
}
