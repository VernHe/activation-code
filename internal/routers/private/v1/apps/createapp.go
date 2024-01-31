package apps

import (
	"configuration-management/global"
	"configuration-management/internal/biz/apps"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type CreateAppRequest struct {
	Name       string `json:"name" form:"name" binding:"required"`
	CardLength int    `json:"card_length" form:"card_length" binding:"min=1"`
	CardPrefix string `json:"card_prefix" form:"card_prefix"`
}

// CreateApp 批量获取app信息
func (handler *Handler) CreateApp(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	err := handler.AppService.CreateApp(apps.CreateAppArgs{
		Name:       req.Name,
		CardLength: req.CardLength,
		CardPrefix: req.CardPrefix,
	})
	if err != nil {
		global.Logger.Error("create app failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
