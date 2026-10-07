package v2

import (
	"errors"

	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/1Panel-dev/1Panel/agent/utils/mirrors"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) ListMirrors(c *gin.Context) {
	data, err := mirrorService.List()
	if err != nil {
		writeMirrorError(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) ApplyMirror(c *gin.Context) {
	var req request.MirrorApplyReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := mirrorService.Apply(req); err != nil {
		writeMirrorError(c, err)
		return
	}
	helper.Success(c)
}

func writeMirrorError(c *gin.Context, err error) {
	var invalid *mirrors.ValidationError
	if errors.As(err, &invalid) {
		helper.BadRequest(c, err)
		return
	}
	helper.InternalServer(c, err)
}
