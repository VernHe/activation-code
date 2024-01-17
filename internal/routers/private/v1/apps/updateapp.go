package apps

import (
	"configuration-management/global"
	"configuration-management/internal/biz/apps"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type UpdateAppRequest struct {
	ID         string `json:"id" form:"id" binding:"required"`
	Name       string `json:"name" form:"name" binding:"required"`
	CardLength int    `json:"card_length" form:"card_length" binding:"min=1,max=32"`
	CardPrefix string `json:"card_prefix" form:"card_prefix"`
}

// UpdateApp 批量获取app信息
func (handler *Handler) UpdateApp(c *gin.Context) {
	var req UpdateAppRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	err := handler.AppService.UpdateApp(apps.UpdateAppArgs{
		ID:         req.ID,
		Name:       req.Name,
		CardLength: req.CardLength,
		CardPrefix: req.CardPrefix,
	})
	if err != nil {
		global.Logger.Error("update app failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
