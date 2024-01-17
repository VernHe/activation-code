package card

import (
	"configuration-management/global"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// GetCardCountByStatus 获取所有状态的激活码数量
func (handler *Handler) GetCardCountByStatus(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var userId string
	if !userInfo.IsRoot() {
		userId = userInfo.UserId
	}

	statusCount, err := handler.CardService.GetCardCountByUserIdAndStatus(userId)
	if err != nil {
		global.Logger.Error("get card failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK(statusCount)
	return
}
