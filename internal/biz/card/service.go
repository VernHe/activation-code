package card

import (
	"time"

	"configuration-management/internal/biz/common"
)

type CreateCardArgs struct {
	UserID string `json:"user_id"` // 用户ID，关联到用户表中的id字段
	IMEI   string `json:"imei"`    // 使用的设备IMEI
	Days   int    `json:"days"`    // 有效天数
}

type CreateCardsArgs struct {
	UserID   string `json:"user_id"`   // 用户ID，关联到用户表中的id字段
	UserName string `json:"user_name"` // 用户ID，关联到用户表中的id字段
	MaxCnt   int    `json:"max_cnt"`   // 最大生成数量
	Minutes  int    `json:"minutes"`   // 有效分钟数
	Hours    int    `json:"hours"`     // 有效小时数
	Days     int    `json:"days"`      // 有效天数
	TimeType string `json:"time_type"` // 时间类型
	Remark   string `json:"remark"`    // 备注信息
	Count    int    `json:"count"`     // 生成数量
	AppID    string `json:"app_id"`    // 应用ID
}

type UpdateCardArgs struct {
	UserId        string `json:"user_id"`         // 用户ID，关联到用户表中的id字段
	CurrentCard   Card   `json:"current_card"`    // 当前卡信息
	ID            string `json:"id"`              // 激活码唯一标识符(UUID)
	Status        int    `json:"status"`          // 状态: 0-未使用, 1-已使用, 2-已锁定, 3-已删除
	Minutes       int    `json:"minutes"`         // 有效分钟数
	Hours         int    `json:"hours"`           // 有效小时数
	Days          int    `json:"days"`            // 有效天数
	TimeType      string `json:"time_type"`       // 时间类型
	IMEI          string `json:"imei"`            // 使用的设备IMEI
	SEID          string `json:"seid"`            // 使用的设备SEID
	Remark        string `json:"remark"`          // 备注信息
	KeepExpiredAt bool   `json:"keep_expired_at"` // 是否保留过期时间
}

type GetCardsArgs struct {
	UserId             string           `json:"user_id"`               // 用户ID，关联到用户表中的id字段
	Values             []string         `json:"values"`                // 激活码值
	Status             []int            `json:"status"`                // 状态: 0-未使用, 1-已使用, 2-已锁定, 3-已删除
	Remark             string           `json:"remark"`                // 备注信息
	NeedPagination     bool             `json:"need_pagination"`       // 是否需要分页
	AppIDs             []string         `json:"app_ids"`               // 应用ID
	TimeType           string           `json:"time_type"`             // 时间类型
	UserName           string           `json:"user_name"`             // 用户ID，关联到用户表中的id字段
	SEID               string           `json:"seid"`                  // 使用的设备SEID
	CreatedAtDateRange common.TimeRange `json:"created_at_date_range"` // 创建时间范围
	UsedAtDateRange    common.TimeRange `json:"used_at_date_range"`    // 使用时间范围
	Page               int              `json:"page"`                  // 页码
	Limit              int              `json:"limit"`                 // 每页数量
}

type GetCardsResult struct {
	Total int    `json:"total"` // 总数
	List  []Card `json:"list"`  // 列表
}

type BatchUpdateStatusArgs struct {
	UserId string   `json:"user_id"`                              // 用户ID，关联到用户表中的id字段
	Values []string `json:"values"`                               // 激活码值
	Status int      `json:"status"`                               // 状态: 0-未使用, 2-已锁定, 3-已删除
	IMEI   string   `json:"imei"`                                 // 使用的设备IMEI
	SEID   string   `json:"seid" gorm:"default:NULL;Column:seid"` // 使用的设备SEID
}

type ActivateCardArgs struct {
	Value string `json:"value"` // 激活码值
	SEID  string `json:"seid"`  // 使用的设备SEID
}

type CheckCardStatusArgs struct {
	Value string `json:"value"` // 激活码值
	SEID  string `json:"seid"`  // 使用的设备SEID
}

type SetCardExpiredAtArgs struct {
	Value        string    `json:"value"`          // 激活码值
	UserId       string    `json:"user_id"`        // 用户ID，关联到用户表中的id字段
	NewExpiredAt time.Time `json:"new_expired_at"` // 新的过期时间
}

type Service interface {
	GetCardByID(id string) (Card, error)
	GetCardByValue(value string) (Card, error)
	GetCardsByUserId(userId string) ([]Card, error)
	GetCards(args GetCardsArgs) (GetCardsResult, error)
	DeleteCardByValue(value string, userId string) error
	CreateCard(args CreateCardArgs) (Card, error)
	CreateCards(args CreateCardsArgs) ([]Card, error)
	UpdateCard(args UpdateCardArgs) error
	DeleteCard(card Card) error
	DeleteCardsByValues(values []string, userId string) error
	BatchUpdateStatus(args BatchUpdateStatusArgs) error
	GetCardCountByUserIdAndStatus(userId string) (map[int]int, error)
	CheckCardStatus(args CheckCardStatusArgs) (bool, error)
	ActivateCard(args ActivateCardArgs) (Card, error)
	SetCardExpiredAt(args SetCardExpiredAtArgs) error
}
