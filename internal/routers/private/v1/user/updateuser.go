package user

import (
	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UpdateUserRequest struct {
	ID           string           `json:"id" binding:"required"`
	MaxCnt       int              `json:"max_cnt" binding:"required"`
	Status       int              `json:"status" binding:"required"`
	Apps         user.Apps        `json:"apps" binding:"required"`
	Permissions  user.Permissions `json:"permissions" binding:"required"`
	Introduction string           `json:"introduction" binding:"required"`
}

func (handler *Handler) UpdateUser(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.WithFields(logger.Fields{
			"req":      req,
			"userInfo": userInfo,
		}).Error("更新用户失败", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 不能自己更新自己
	if req.ID == userInfo.UserId {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("不能自己更新自己"))
		return
	}

	if err := handler.UserService.UpdateUser(user.UpdateUserArgs{
		UpdaterID:    userInfo.UserId,
		ID:           req.ID,
		MaxCnt:       req.MaxCnt,
		Apps:         req.Apps,
		Permissions:  req.Permissions,
		Introduction: req.Introduction,
		Status:       req.Status,
	}); err != nil {
		global.Logger.WithFields(logger.Fields{
			"req":    req,
			"userId": userInfo.UserId,
		}).Error("update user failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
}
