package card

import (
	"configuration-management/pkg/app"

	"github.com/gin-gonic/gin"
)

type CreateCardRequest struct {
	Days int `json:"days" binding:"required" form:"days"`
}

// CreateCard 创建激活码
// 目前未使用
func (handler *Handler) CreateCard(c *gin.Context) {
	app.NewResponse(c).ResponseOK("ok")
	return

	//var req CreateCardRequest
	//if err := c.ShouldBind(&req); err != nil {
	//	global.Logger.WithFields(logger.Fields{
	//		"request_body": c.Request.Body,
	//	}).Error("invalid params", err)
	//	app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	//	return
	//}
	//
	//newCard, err := handler.CardService.CreateCard(card.CreateCardArgs{
	//	UserID: "fakeId",
	//	IMEI:   "fakeImei",
	//	Days:   req.Days,
	//})
	//if err != nil {
	//	global.Logger.WithFields(logger.Fields{
	//		"error": err,
	//	}).Error("create card failed")
	//	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	//	return
	//}
	//
	//app.NewResponse(c).ResponseOK(newCard)
	//return
}
