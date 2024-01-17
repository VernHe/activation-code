package card

import (
	"strings"

	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type BatchUpdateStatusRequest struct {
	Values string `json:"values" binding:"required"`
	Status int    `json:"status" binding:"required"`
	SEID   string `json:"seid"`
}

// BatchUpdateStatus 批量更新卡信息
func (handler *Handler) BatchUpdateStatus(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req BatchUpdateStatusRequest
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

	var userId string
	if !userInfo.IsRoot() {
		currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"userInfo": userInfo,
			}).Error("用户不存在", err)
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission.WithDetails(err.Error()))
			return
		}

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

		userId = userInfo.UserId
	}

	if err := handler.CardService.BatchUpdateStatus(card.BatchUpdateStatusArgs{
		UserId: userId,
		Values: strings.Split(req.Values, ","),
		Status: req.Status,
		SEID:   req.SEID,
	}); err != nil {
		global.Logger.Error("batch update status failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
}
