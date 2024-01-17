package card

import (
	"time"

	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type SetCardExpiredAtRequest struct {
	Value        string    `json:"value" binding:"required"`
	NewExpiredAt time.Time `json:"new_expired_at" binding:"required"`
}

// SetCardExpiredAt 设置激活码过期时间
func (handler *Handler) SetCardExpiredAt(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req SetCardExpiredAtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 检查用户是否有权限
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
		// 检查用户是否有查询权限
		if !currentUser.HasPermission(permissions.UPDATE) {
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
			return
		}
		userId = currentUser.ID
	}

	if err := handler.CardService.SetCardExpiredAt(card.SetCardExpiredAtArgs{
		Value:        req.Value,
		UserId:       userId,
		NewExpiredAt: req.NewExpiredAt,
	}); err != nil {
		global.Logger.Error("update card failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
