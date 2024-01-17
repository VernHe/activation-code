package app

import (
	"net/http"

	"configuration-management/pkg/errcode"

	"github.com/gin-gonic/gin"
)

const (
	UserInfoKey = "userInfo"
)

var (
	EmptyResponseContent = ResponseContent{StatusCode: http.StatusOK, Data: map[string]any{}}
)

type Response struct {
	Ctx *gin.Context
}

type ResponseContent struct {
	StatusCode int `json:"code"`
	Data       any `json:"data"`
}

type Pager struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ResponseOK(data ...any) {
	if len(data) > 0 {
		r.Ctx.JSON(http.StatusOK, ResponseContent{
			StatusCode: http.StatusOK,
			Data:       data[0],
		})
		return
	}
	r.Ctx.JSON(http.StatusOK, EmptyResponseContent)
}

func (r *Response) ToResponse(content ResponseContent) {
	r.Ctx.JSON(content.StatusCode, content)
}

func (r *Response) ToResponseList(list interface{}, total int) {
	content := ResponseContent{
		StatusCode: http.StatusOK,
		Data: map[string]any{
			"items": list,
			"pager": Pager{
				Page:     GetPage(r.Ctx),
				PageSize: GetPageSize(r.Ctx),
				Total:    total,
			},
		},
	}
	r.Ctx.JSON(http.StatusOK, content)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}

func GetUserInfoFromContext(c *gin.Context) UserInfo {
	return c.MustGet(UserInfoKey).(UserInfo)
}
