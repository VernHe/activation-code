package card

import (
	"configuration-management/global"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type GetCardRequest struct {
	Value string `json:"value" binding:"required"`
}

// GetCardByValue 根据激活码获取卡信息
func (handler *Handler) GetCardByValue(c *gin.Context) {
	value := c.Param("value")
	if value == "" {
		global.Logger.WithFields(logger.Fields{
			"value": value,
		}).Error("查询空的激活码")
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("value", value))
		return
	}

	card, err := handler.CardService.GetCardByValue(value)
	if err != nil {
		global.Logger.Error("get card failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK(card)
	return
}
