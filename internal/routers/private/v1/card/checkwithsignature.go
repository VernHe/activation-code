package card

import (
	"encoding/json"
	"fmt"

	"configuration-management/internal/biz/card"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/utils/security"

	"github.com/gin-gonic/gin"
)

type CheckDecryptedData struct {
	Value string `json:"value"` // 激活码值
	SEID  string `json:"seid"`  // 使用的设备SEID
}

type CheckRequestBody struct {
	EncryptedData string `json:"data"`
	Signature     string `json:"signature"`
	Timestamp     string `json:"timestamp"`
}

type CheckResponseBody struct {
	Result    string `json:"rs"`
	ExtraData string `json:"x"`
	Signature string `json:"s"`
}

// CheckWithSignature 检查激活码状态，带签名
func (handler *Handler) CheckWithSignature(c *gin.Context) {
	var req CheckRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 校验时间戳
	if !security.IsValidTimestamp(req.Timestamp) {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("参数错误"))
		return
	}

	// 校验签名
	if !security.IsValidSignature(req.Signature, req.Timestamp, req.EncryptedData) {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("参数错误"))
		return
	}

	// 解密
	decrypted, err := security.GetAESDecrypted(req.EncryptedData)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	// 解密后的数据
	var data CheckDecryptedData
	if err := json.Unmarshal(decrypted, &data); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 业务逻辑
	activated, err := handler.CardService.CheckCardStatus(card.CheckCardStatusArgs{
		Value: data.Value,
		SEID:  data.SEID,
	})
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	var response CheckResponseBody
	// Result
	encryptedResult, err := security.GetAESEncrypted(fmt.Sprintf("%t", activated))
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.Result = encryptedResult
	// ExtraData，保存暗号
	if activated {
		response.ExtraData = security.GenerateCipherText(data.Value, false)
	} else {
		response.ExtraData = security.GenerateCipherText(data.Value, true)
	}

	// 加载 RSA 私钥，用于签名，签名内容为 Result + ExtraData
	signature, err := security.GetSignature(fmt.Sprintf("%s%s", response.Result, response.ExtraData))
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.Signature = signature

	// 加密返回的数据
	app.NewResponse(c).ResponseOK(response)
}
