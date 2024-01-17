package user

import (
	"net/http"

	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type GetUserInfoResponse struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

func (handler *Handler) GetUserInfo(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	targetUser, err := handler.UserService.GetUserInfo(user.GetUserInfoArgs{UserId: userInfo.UserId})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"userId": userInfo.UserId,
		}).Error("login failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	var response GetUserInfoResponse
	response.Name = targetUser.Username
	response.Roles = targetUser.Roles
	response.Avatar = targetUser.Avatar
	response.Introduction = targetUser.Introduction

	app.NewResponse(c).ToResponse(app.ResponseContent{
		StatusCode: http.StatusOK,
		Data:       response,
	})
	return
}
