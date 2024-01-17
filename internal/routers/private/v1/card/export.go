package card

import (
	"strings"

	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type ExportRequest struct {
	Value  string `form:"value"`
	Status int    `form:"status"`
	Remark string `form:"remark"`
	Page   int    `form:"page" binding:"required,min=1"`
	Limit  int    `form:"limit" binding:"required,min=1,max=50"`
}

// Export 导出卡信息的 JSON
func (handler *Handler) Export(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req GetCardsRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	var values []string
	if req.Value != "" {
		values = strings.Split(req.Value, ",")
	}
	var status []int
	if req.Status != 0 {
		status = append(status, req.Status)
	}

	result, err := handler.CardService.GetCards(card.GetCardsArgs{
		UserId:         userInfo.UserId,
		Values:         values,
		Status:         status,
		Remark:         req.Remark,
		NeedPagination: false,
	})
	if err != nil {
		global.Logger.Error("get cards failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	appOptions, err := handler.AppService.QueryAppOptions()
	if err != nil {
		global.Logger.Error("get app options failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ToResponseList(card.BatchToView(result.List, appOptions), result.Total)
	return
}
