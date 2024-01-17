package card

import (
	"strconv"
	"strings"
	"time"

	"configuration-management/global"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/common"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type BatchQueryRequest struct {
	Values             string       `form:"values"`
	Status             string       `form:"status"`
	Remark             string       `form:"remark"`
	AppIds             string       `form:"app_ids"`
	UserName           string       `form:"user_name"`
	CreatedAtDateRange [2]time.Time `form:"created_at_date_range[]"`
	UsedAtDateRange    [2]time.Time `form:"used_at_date_range[]"`
}

// BatchQuery 批量查询卡信息
func (handler *Handler) BatchQuery(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req BatchQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
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
		// 检查用户是否有查询权限
		if !currentUser.HasPermission(permissions.QUERY) {
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
			return
		}

		userId = userInfo.UserId
	}
	var appIds, values []string
	if req.AppIds != "" {
		appIds = strings.Split(req.AppIds, ",")
	}
	if req.Values != "" {
		values = strings.Split(req.Values, ",")
	}
	// string split 后再转 int slice
	var status []int
	if req.Status != "" {
		for _, s := range strings.Split(req.Status, ",") {
			// string 转 int
			i, err := strconv.Atoi(s)
			if err != nil {
				global.Logger.WithFields(logger.Fields{
					"status": req.Status,
				}).Error("invalid status", err)
				app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
				return
			}
			status = append(status, i)
		}
	}

	result, err := handler.CardService.GetCards(card.GetCardsArgs{
		UserId:         userId,
		AppIDs:         appIds,
		UserName:       req.UserName,
		Values:         values,
		Status:         status,
		Remark:         req.Remark,
		NeedPagination: false,
		CreatedAtDateRange: common.TimeRange{
			StartTime: req.CreatedAtDateRange[0],
			EndTime:   req.CreatedAtDateRange[1],
		},
		UsedAtDateRange: common.TimeRange{
			StartTime: req.UsedAtDateRange[0],
			EndTime:   req.UsedAtDateRange[1],
		},
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
