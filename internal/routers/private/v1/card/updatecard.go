package card

import (
	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UpdateCardRequest struct {
	ID       string `json:"id" binding:"required"`
	Status   int    `json:"status"`
	AppId    string `json:"app_id"`
	SEID     string `json:"seid"`
	TimeType string `json:"time_type"`
	Minutes  int    `json:"minutes"`
	Horus    int    `json:"hours"`
	Days     int    `json:"days"`
	Remark   string `json:"remark"`
}

// UpdateCard 更新卡信息
func (handler *Handler) UpdateCard(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 读取 request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = c.GetRawData()
		}

		global.Logger.WithFields(logger.Fields{
			"req_body": string(bodyBytes),
		}).Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 检查用户是否有权限
	var userId string
	isRoot := userInfo.IsRoot()
	currentCard, err := handler.CardService.GetCardByID(req.ID)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"cardId": req.ID,
		}).Error("查询激活码失败", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}
	if !isRoot {
		currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"userInfo": userInfo,
			}).Error("用户不存在", err)
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission.WithDetails(err.Error()))
			return
		}
		// 检查用户是否有查询权限
		if !currentUser.HasPermission(permissions.UPDATE) {
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
			return
		}

		// 检查是否有更新状态
		if req.Status != 0 && req.Status != currentCard.Status {
			var subPermission string
			switch req.Status {
			case card.StatusUsed:
				subPermission = permissions.UPDATE_USED
			case card.StatusUnused:
				subPermission = permissions.UPDATE_UNUSED
			case card.StatusLocked:
				subPermission = permissions.UPDATE_LOCKED
			case card.StatusDeleted:
				subPermission = permissions.UPDATE_DELETED
			default:
				global.Logger.WithFields(logger.Fields{
					"status": req.Status,
				}).Error("invalid status")
				app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("invalid status"))
				return
			}
			// 检查用户是否有更新权限
			if !currentUser.HasPermission(subPermission) {
				app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
				return
			}
		}
	}

	// 检查时间类型和时间
	//if !biz.IsValidTimeType(req.TimeType, req.Days, req.Minutes) {
	//	global.Logger.WithFields(logger.Fields{
	//		"req":  req,
	//		"user": userInfo,
	//	}).Error("时间类型和时间不匹配")
	//	app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("时间类型和时间不匹配"))
	//	return
	//}

	if err := handler.CardService.UpdateCard(card.UpdateCardArgs{
		UserId:      userId,
		CurrentCard: currentCard,
		ID:          req.ID,
		Status:      req.Status,
		//Minutes:     req.Minutes,
		Hours:    req.Horus,
		Days:     req.Days,
		TimeType: req.TimeType,
		SEID:     req.SEID,
		Remark:   req.Remark,
	}); err != nil {
		global.Logger.Error("update card failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
