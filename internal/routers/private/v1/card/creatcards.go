package card

import (
	"net/http"

	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type CreateCardsRequest struct {
	Minutes  int    `json:"minutes"`
	Hours    int    `json:"hours"`
	Days     int    `json:"days"`
	TimeType string `json:"time_type"`
	Count    int    `json:"count" binding:"required"`
	AppID    string `json:"app_id" binding:"required"`
	Remark   string `json:"remark"`
}

// CreateCards 批量创建激活码
func (handler *Handler) CreateCards(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req CreateCardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.WithFields(logger.Fields{
			"request_body": c.Request.Body,
		}).Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 权限校验
	currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"error": err,
		}).Error("[CreateCards] 查询用户失败")
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
		return
	}
	if !userInfo.IsRoot() {
		// 检查用户是否有查询权限
		if !currentUser.HasPermission(permissions.CREATE) {
			global.Logger.WithFields(logger.Fields{
				"error": err,
			}).Error("用户权限异常，无法创建激活码")
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
			return
		}
	}

	// 对时间类型和时间进行校验
	//if !biz.IsValidTimeType(req.TimeType, req.Days, req.Minutes) {
	//	global.Logger.WithFields(logger.Fields{
	//		"req":  req,
	//		"user": userInfo,
	//	}).Error("时间类型和时间不匹配")
	//	app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("时间类型和时间不匹配"))
	//	return
	//}

	// 创建
	newCards, err := handler.CardService.CreateCards(card.CreateCardsArgs{
		UserID:   userInfo.UserId,
		UserName: userInfo.Username,
		MaxCnt:   currentUser.MaxCnt,
		//Minutes:  req.Minutes,
		Hours:    req.Hours,
		Days:     req.Days,
		TimeType: req.TimeType,
		Remark:   req.Remark,
		Count:    req.Count,
		AppID:    req.AppID,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"error": err,
		}).Error("create card failed")
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
		return
	}

	app.NewResponse(c).ToResponse(app.ResponseContent{
		StatusCode: http.StatusOK,
		Data:       newCards,
	})
	return
}
