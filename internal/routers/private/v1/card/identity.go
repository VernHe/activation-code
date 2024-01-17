package card

import (
	"errors"
	"fmt"
	"time"

	"configuration-management/internal/biz/card"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type IdentityRequest struct {
	Value string `json:"value"` // 激活码值
	SEID  string `json:"seid"`  // 使用的设备SEID
}

type IdentityResponse struct {
	Card card.Card `json:"card"`
}

func (handler *Handler) Identity(c *gin.Context) {
	var req IdentityRequest
	var resp IdentityResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 校验参数
	if req.Value == "" || req.SEID == "" {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("参数错误"))
		return
	}

	// 检查激活码是否存在
	code, err := handler.CardService.GetCardByValue(req.Value)
	if err != nil {
		if errors.Is(err, errcode.NotFound) {
			app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
			return
		}
		fmt.Println(err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	// 检查激活码的状态是否合法
	if code.Status == card.StatusLocked || code.Status == card.StatusDeleted {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("激活码不可用"))
		return
	}

	// 激活码第一次使用，进行激活
	if code.Status == card.StatusUnused {
		activatedCode, err := handler.CardService.ActivateCard(card.ActivateCardArgs{
			Value: req.Value,
			SEID:  req.SEID,
		})
		if err != nil {
			app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
			return
		}
		resp.Card = activatedCode
		app.NewResponse(c).ResponseOK(resp)
		return
	}

	// 检查是否过期
	if code.ExpiredAt != nil && time.Now().After(*code.ExpiredAt) {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("激活码已过期"))
		return
	}

	// 激活码已经被使用，检查设备是否匹配
	if code.Status == card.StatusUsed && code.SEID != req.SEID {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("设备不匹配"))
		return
	}

	resp.Card = code
	app.NewResponse(c).ResponseOK(resp)
}
