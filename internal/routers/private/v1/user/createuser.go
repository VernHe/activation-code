package user

import (
	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username     string           `json:"username" binding:"required"`
	Password     string           `json:"password" binding:"required"`
	MaxCnt       int              `json:"max_cnt" binding:"required"`
	Permissions  user.Permissions `json:"permissions"`
	Apps         user.Apps        `json:"apps"`
	Introduction string           `json:"introduction"`
}

func (handler *Handler) CreateUser(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.WithFields(logger.Fields{
			"req":      req,
			"userInfo": userInfo,
		}).Error("create user failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if err := handler.UserService.CreateUser(user.CreateUserArgs{
		CreatorID:    userInfo.UserId,
		Username:     req.Username,
		Password:     req.Password,
		MaxCnt:       req.MaxCnt,
		Apps:         req.Apps,
		Permissions:  req.Permissions,
		Introduction: req.Introduction,
	}); err != nil {
		global.Logger.WithFields(logger.Fields{
			"req":    req,
			"userId": userInfo.UserId,
		}).Error("create user failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.DuplicateKey.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
}
