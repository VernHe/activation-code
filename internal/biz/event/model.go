package event

import "time"

// Event 结构体对应 Event 表
type Event struct {
	ID           int       `json:"id"`             // 事件唯一标识符
	UserID       int       `json:"user_id"`        // 事件所属用户的ID，关联到用户表中的id字段
	RemindUserID int       `json:"remind_user_id"` // 提醒的用户ID（如果有的话），关联到用户表中的id字段
	Name         string    `json:"name"`           // 事件名称
	Detail       string    `json:"detail"`         // 事件详细描述
	Read         bool      `json:"read"`           // 是否已读
	CreatedAt    time.Time `json:"created_at"`     // 事件创建时间
}
