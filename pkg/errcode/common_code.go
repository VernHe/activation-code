package errcode

var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(10000000, "服务内部错误")
	InvalidParams             = NewError(10000001, "入参错误")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败，找不到对应的 AppKey 和 AppSecret")
	UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token 错误")
	UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token 超时")
	UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token 生成失败")
	TooManyRequests           = NewError(10000007, "请求过多")
	DuplicateKey              = NewError(10000008, "数据已存在")
	NoPermission              = NewError(10000009, "没有权限")

	CardNotFound       = NewError(20010000, "激活码找不到")
	CardNotAvailable   = NewError(20010001, "激活码不可用")
	DeviceNotAvailable = NewError(20010002, "设备不可用")
	CardExpired        = NewError(20010003, "激活码已过期")
)
