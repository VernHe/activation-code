package card

import (
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

type GetCardsRequest struct {
	Value              string       `form:"value"`
	Status             int          `form:"status"`
	Remark             string       `form:"remark"`
	AppId              string       `form:"app_id"`
	UserName           string       `form:"user_name"`
	TimeType           string       `form:"time_type"`
	SEID               string       `form:"seid"`
	CreatedAtDateRange [2]time.Time `form:"created_at_date_range[]"`
	UsedAtDateRange    [2]time.Time `form:"used_at_date_range[]"`
	UsedAtDate         string       `form:"used_at"`
	NeedPagination     bool         `form:"need_pagination"`
	Page               int          `form:"page"`
	Limit              int          `form:"limit"`
}

// GetCards 批量获取卡信息
func (handler *Handler) GetCards(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req GetCardsRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	var appIds []string
	if req.AppId != "" {
		appIds = append(appIds, strings.Split(req.AppId, ",")...)
	}
	var values []string
	if req.Value != "" {
		values = append(values, strings.Split(req.Value, ",")...)
	}
	var status []int
	if req.Status != 0 {
		status = append(status, req.Status)
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

		// 检验 appIds 是否在用户的权限范围内
		if len(appIds) > 0 {
			for _, appId := range appIds {
				// 检查是否有使用该 app
				if !currentUser.HasApp(appId) {
					global.Logger.WithFields(logger.Fields{
						"userInfo": userInfo,
						"appId":    appId,
					}).Error("用户没有使用该 app 的权限")
					app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
					return
				}
			}
		} else {
			// 没有指定 appIds，则使用用户有权限的 appIds
			appIds = currentUser.Apps
		}

		userId = userInfo.UserId
	}

	// 默认分页值
	if req.NeedPagination {
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.Limit <= 0 {
			req.Limit = 10
		}
	}

	// 时间范围
	var usedAtDateRange common.TimeRange
	if req.UsedAtDate != "" {
		usedAtDate, err := time.Parse("2006-01-02", req.UsedAtDate)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"used_at_date": req.UsedAtDate,
			}).Error("invalid used_at_date", err)
			app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
			return
		}

		// 范围是当天的 00:00:00 ~ 23:59:59
		usedAtDateRange = common.TimeRange{
			StartTime: usedAtDate,
			EndTime:   usedAtDate.Add(24 * time.Hour).Add(-1 * time.Second),
		}
	}
	if (usedAtDateRange.StartTime.IsZero() || usedAtDateRange.EndTime.IsZero()) && req.UsedAtDateRange[0].IsZero() && req.UsedAtDateRange[1].IsZero() {
		usedAtDateRange.StartTime = req.UsedAtDateRange[0]
		usedAtDateRange.EndTime = req.UsedAtDateRange[1]
	}

	result, err := handler.CardService.GetCards(card.GetCardsArgs{
		UserId:         userId,
		AppIDs:         appIds,
		TimeType:       req.TimeType,
		UserName:       req.UserName,
		Values:         values,
		Status:         status,
		Remark:         req.Remark,
		NeedPagination: req.NeedPagination,
		SEID:           req.SEID,
		CreatedAtDateRange: common.TimeRange{
			StartTime: req.CreatedAtDateRange[0].UTC(),
			EndTime:   req.CreatedAtDateRange[1].UTC(),
		},
		UsedAtDateRange: usedAtDateRange,
		Page:            req.Page,
		Limit:           req.Limit,
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
