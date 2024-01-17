package card

import (
	"time"

	"configuration-management/internal/biz/apps"
)

// Card 结构体对应 Card 表
type Card struct {
	ID        string     `json:"id"`                                                      // 激活码唯一标识符(UUID)
	AppID     string     `json:"app_id"`                                                  // 应用ID
	Status    int        `json:"status"`                                                  // 状态: 1-未使用, 2-已使用, 3-已锁定, 4-已删除
	UserID    string     `json:"user_id"`                                                 // 用户ID，关联到用户表中的id字段
	UserName  string     `json:"user_name"`                                               // 创建激活码的用户名
	Days      int        `json:"days"`                                                    // 有效天数
	Minutes   int        `json:"minutes"`                                                 // 有效分钟数
	ExpiredAt *time.Time `json:"expired_at" gorm:"default:NULL type:timestamp"`           // 过期时间
	TimeType  string     `json:"time_type"`                                               // 时间类型 hourly-小时, daily-天数, weekly-周数, monthly-月数, yearly-年数
	Value     string     `json:"value"`                                                   // 激活码值
	Used      bool       `json:"used"`                                                    // 是否使用过
	SEID      string     `json:"seid" gorm:"default:NULL;Column:seid"`                    // 使用的设备SEID
	UsedAt    *time.Time `json:"used_at" gorm:"default:NULL type:timestamp"`              // 使用时间
	LockedAt  *time.Time `json:"locked_at" gorm:"default:NULL type:timestamp"`            // 锁定时间
	Remark    string     `json:"remark"`                                                  // 备注信息
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL type:timestamp"` // 删除时间
	CreatedAt *time.Time `json:"created_at"`                                              // 生成时间
}

// CardView 返回给前端的激活码信息
type CardView struct {
	ID        string `json:"id"`                   // 激活码唯一标识符(UUID)
	AppID     string `json:"app_id"`               // 应用ID
	AppName   string `json:"app_name"`             // 应用名称
	Status    int    `json:"status"`               // 状态: 0-未使用, 1-已使用, 2-已锁定, 3-已删除
	UserID    string `json:"user_id"`              // 用户ID，关联到用户表中的id字段
	UserName  string `json:"user_name"`            // 创建激活码的用户名
	Days      int    `json:"days"`                 // 有效天数
	Minutes   int    `json:"minutes"`              // 有效分钟数
	ExpiredAt string `json:"expired_at"`           // 过期时间
	TimeType  string `json:"time_type"`            // 时间类型 hourly-小时, daily-天数, weekly-周数, monthly-月数, yearly-年数
	Value     string `json:"value"`                // 激活码值
	Used      bool   `json:"used"`                 // 是否使用过
	SEID      string `json:"seid"`                 // 使用的设备SEID
	UsedAt    string `json:"used_at"`              // 使用时间
	LockedAt  string `json:"locked_at"`            // 锁定时间
	Remark    string `json:"remark"`               // 备注信息
	DeletedAt string `json:"deleted_at,omitempty"` // 删除时间
	CreatedAt string `json:"created_at"`           // 生成时间
}

// TableName 指定 Card 结构体对应的表名
func (card *Card) TableName() string {
	return "card"
}

func (card *Card) ToView(appOptions []apps.AppOption) CardView {
	// 避免空指针
	var expiredAt, usedAt, locakedAt, deletedAt, createdAt string
	if card.ExpiredAt != nil {
		expiredAt = card.ExpiredAt.Format("2006-01-02 15:04:05")
	}
	if card.UsedAt != nil {
		usedAt = card.UsedAt.Format("2006-01-02 15:04:05")
	}
	if card.LockedAt != nil {
		locakedAt = card.LockedAt.Format("2006-01-02 15:04:05")
	}
	if card.DeletedAt != nil {
		deletedAt = card.DeletedAt.Format("2006-01-02 15:04:05")
	}
	if card.CreatedAt != nil {
		createdAt = card.CreatedAt.Format("2006-01-02 15:04:05")
	}

	appName := ""
	for _, option := range appOptions {
		if option.ID == card.AppID {
			appName = option.Name
			break
		}
	}

	return CardView{
		ID:        card.ID,
		AppID:     card.AppID,
		AppName:   appName,
		Status:    card.Status,
		UserID:    card.UserID,
		UserName:  card.UserName,
		Days:      card.Days,
		Minutes:   card.Minutes,
		ExpiredAt: expiredAt,
		TimeType:  card.TimeType,
		Value:     card.Value,
		Used:      card.Used,
		SEID:      card.SEID,
		UsedAt:    usedAt,
		LockedAt:  locakedAt,
		Remark:    card.Remark,
		DeletedAt: deletedAt,
		CreatedAt: createdAt,
	}

}

func BatchToView(cards []Card, appOptions []apps.AppOption) []CardView {
	var views []CardView
	for _, card := range cards {
		views = append(views, card.ToView(appOptions))
	}
	return views
}
