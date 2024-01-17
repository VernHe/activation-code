package card

import (
	"encoding/json"
	"fmt"
	"time"

	"configuration-management/global"
	"configuration-management/internal/biz/activationattempt"
	"configuration-management/internal/biz/card"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"
	"configuration-management/utils/security"

	"github.com/gin-gonic/gin"
)

type ActivateDecryptedData struct {
	Value string `json:"value"` // 激活码值
	SEID  string `json:"seid"`  // 使用的设备SEID
}

type ActivateRequestBody struct {
	EncryptedData string `json:"data"`
	Signature     string `json:"signature"`
	Timestamp     string `json:"timestamp"`
}

type ActivateResponseBody struct {
	Result    string `json:"rs"`
	ExtraData string `json:"x"`
	Signature string `json:"s"`
}

// Activate 检查激活码状态
func (handler *Handler) Activate(c *gin.Context) {
	var req ActivateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 校验时间戳
	if !security.IsValidTimestamp(req.Timestamp) {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("invalid timestamp"))
		return
	}

	// 校验签名
	if !security.IsValidSignature(req.Signature, req.Timestamp, req.EncryptedData) {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("invalid signature"))
		return
	}

	// 解密
	decrypted, err := security.GetAESDecrypted(req.EncryptedData)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	// 解密后的数据
	var data ActivateDecryptedData
	if err := json.Unmarshal(decrypted, &data); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	requestData, err := json.Marshal(data)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	// 业务逻辑
	_, err = handler.CardService.ActivateCard(card.ActivateCardArgs{
		Value: data.Value,
		SEID:  data.SEID,
	})
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		if err = handler.ActivationAttempt.CreateActivationAttempt(&activationattempt.ActivationAttempt{
			CardValue:    data.Value,
			Success:      false,
			ErrorMessage: err.Error(),
			RequestData:  string(requestData),
			ResponseData: err.Error(),
			CreatedAt:    time.Now(),
		}); err != nil {
			global.Logger.WithFields(logger.Fields{
				"card_value": data.Value,
				"request":    string(requestData),
			}).Error("create activation attempt failed", err)
			return
		}
		return
	}

	if err = handler.ActivationAttempt.CreateActivationAttempt(&activationattempt.ActivationAttempt{
		CardValue:    data.Value,
		Success:      true,
		ErrorMessage: "",
		RequestData:  string(requestData),
		ResponseData: "",
		CreatedAt:    time.Now(),
	}); err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	var response ActivateResponseBody
	encryptedResult, err := security.GetAESEncrypted(data.Value)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.Result = encryptedResult

	// 检查激活次数
	totalCount, err := handler.ActivationAttempt.GetActivationAttemptCountByCardValue(data.Value)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	// 检查一小时内激活次数
	hourCount, err := handler.ActivationAttempt.GetActivationAttemptCountByCardValueInHour(data.Value)
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	// ExtraData，暗号
	response.ExtraData = security.GenerateCipherText(data.Value, hourCount > 3 || totalCount > 5)

	// 加载 RSA 私钥，用于签名，签名内容为 Status + ExtraData
	signature, err := security.GetSignature(fmt.Sprintf("%s%s", response.Result, response.ExtraData))
	if err != nil {
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.Signature = signature

	app.NewResponse(c).ResponseOK(response)
}
