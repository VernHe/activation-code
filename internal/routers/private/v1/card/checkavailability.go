package card

import (
	"configuration-management/global"
	card2 "configuration-management/internal/biz/card"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type CheckAvailabilityRequest struct {
	CardValue string `json:"card_value" binding:"required"`
}

type CheckAvailabilityResponse struct {
	IsAvailable bool `json:"is_available"`
}

func (handler *Handler) CheckAvailability(c *gin.Context) {
	var req CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	isAvailable, err := handler.CardService.CheckAvailability(card2.CheckAvailabilityArgs{CardValue: req.CardValue})
	if err != nil {
		global.Logger.Error("get card failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	rsp := CheckAvailabilityResponse{
		IsAvailable: isAvailable,
	}
	app.NewResponse(c).ResponseOK(rsp)
	return
}
