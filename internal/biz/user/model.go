package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Apps []string

func (a Apps) ToJsonArray() json.RawMessage {
	// 存到数据库是: ["aaa", "bbb", "ccc"]，标准对 json 格式，每个元素需要有双引号
	// 读取时是: [aaa,bbb,ccc]，标准对 json 格式，每个元素不需要有双引号
	var apps []string
	for _, app := range a {
		apps = append(apps, app)
	}
	appsJson, _ := json.Marshal(apps)
	return appsJson
}

// Scan 实现了 sql.Scanner 接口
func (a *Apps) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case json.RawMessage:
		if len(v) == 0 {
			*a = nil
			return nil
		}
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return errors.New("不支持的 Scan 操作，将 driver.Value 类型存储到 *Apps 类型中")
	}
}

// Value 实现了 driver.Valuer 接口
func (a *Apps) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Roles 是处理数据库中的 JSON 数组的自定义类型
type Roles []string

func (r Roles) ToJsonArray() json.RawMessage {
	// 存到数据库是: ["aaa", "bbb", "ccc"]，标准对 json 格式，每个元素需要有双引号
	// 读取时是: [aaa,bbb,ccc]，标准对 json 格式，每个元素不需要有双引号
	var roles []string
	for _, role := range r {
		roles = append(roles, role)
	}
	rolesJson, _ := json.Marshal(roles)
	return rolesJson
}

// Scan 实现了 sql.Scanner 接口
func (r *Roles) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case json.RawMessage:
		if len(v) == 0 {
			*r = nil
			return nil
		}
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return errors.New("不支持的 Scan 操作，将 driver.Value 类型存储到 *Roles 类型中")
	}
}

// Value 实现了 driver.Valuer 接口
func (r *Roles) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Permissions 是处理数据库中的 JSON 数组的自定义类型
type Permissions []string

func (p Permissions) ToJsonArray() json.RawMessage {
	// 存到数据库是: ["aaa", "bbb", "ccc"]，标准对 json 格式，每个元素需要有双引号
	// 读取时是: [aaa,bbb,ccc]，标准对 json 格式，每个元素不需要有双引号
	var permissions []string
	for _, permission := range p {
		permissions = append(permissions, permission)
	}
	permissionsJson, _ := json.Marshal(permissions)
	return permissionsJson
}

// Scan 实现了 sql.Scanner 接口
func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case json.RawMessage:
		if len(v) == 0 {
			*p = nil
			return nil
		}
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return errors.New("不支持的 Scan 操作，将 driver.Value 类型存储到 *Permissions 类型中")
	}
}

// Value 实现了 driver.Valuer 接口
func (p *Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type Status int

func (s Status) String() string {
	switch s {
	case StatusNormal:
		return "正常"
	case StatusBanned:
		return "停封"
	default:
		return "未知"
	}
}

func NewStatus(status int) Status {
	switch status {
	case 1:
		return StatusNormal
	case -1:
		return StatusBanned
	default:
		return StatusUnknown
	}
}

type DBStruct struct {
	ID           string          `json:"id"`                           // 用户唯一标识符(UUID)
	Username     string          `json:"username"`                     // 用户名
	Password     string          `json:"password"`                     // 密码
	Status       Status          `json:"status"`                       // 状态: -1-停封, 1-正常
	Ancestry     string          `json:"ancestry"`                     // 祖先用户ID（如果有的话）
	TotalCnt     int             `json:"total_cnt"`                    // 总激活码数量
	UsedCnt      int             `json:"used_cnt"`                     // 已使用激活码数量
	NousedCnt    int             `json:"noused_cnt"`                   // 未使用激活码数量
	DeletedCnt   int             `json:"deleted_cnt"`                  // 删除的激活码数量
	LockedCnt    int             `json:"locked_cnt"`                   // 锁定的激活码数量
	MaxCnt       int             `json:"max_cnt"`                      // 最大激活码数量
	CreatedAt    time.Time       `json:"created_at"`                   // 用户创建时间
	Roles        json.RawMessage `json:"roles" gorm:"type:json"`       // 用户角色
	Introduction string          `json:"introduction"`                 // 介绍
	Avatar       string          `json:"avatar"`                       // 头像
	Permissions  json.RawMessage `json:"permissions" gorm:"type:json"` // 用户权限
	Apps         json.RawMessage `json:"apps" gorm:"type:json"`        // 用户有权限的应用
}

func (s *DBStruct) TableName() string {
	return "user"
}

func (s *DBStruct) ToModel() (User, error) {
	user := User{
		ID:           s.ID,
		Username:     s.Username,
		Password:     s.Password,
		Status:       s.Status,
		Ancestry:     s.Ancestry,
		TotalCnt:     s.TotalCnt,
		UsedCnt:      s.UsedCnt,
		NousedCnt:    s.NousedCnt,
		DeletedCnt:   s.DeletedCnt,
		LockedCnt:    s.LockedCnt,
		MaxCnt:       s.MaxCnt,
		CreatedAt:    s.CreatedAt,
		Introduction: s.Introduction,
		Avatar:       s.Avatar,
	}
	if err := user.Roles.Scan(s.Roles); err != nil {
		return User{}, err
	}
	if err := user.Permissions.Scan(s.Permissions); err != nil {
		return User{}, err
	}
	if err := user.Apps.Scan(s.Apps); err != nil {
		return User{}, err
	}

	return user, nil
}

func BatchToModel(dbStructs []DBStruct) ([]User, error) {
	var users []User
	for _, dbStruct := range dbStructs {
		user, err := dbStruct.ToModel()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// User 结构体对应 User 表
type User struct {
	ID           string      `json:"id" structs:"id"`                     // 用户唯一标识符(UUID)
	Username     string      `json:"username" structs:"username"`         // 用户名
	Password     string      `json:"password" structs:"password"`         // 密码
	Status       Status      `json:"status" structs:"status"`             // 状态: 0-停封, 1-正常
	Ancestry     string      `json:"ancestry" structs:"ancestry"`         // 祖先用户ID（如果有的话）
	TotalCnt     int         `json:"total_cnt" structs:"total_cnt"`       // 总激活码数量
	UsedCnt      int         `json:"used_cnt" structs:"used_cnt"`         // 已使用激活码数量
	NousedCnt    int         `json:"noused_cnt" structs:"noused_cnt"`     // 未使用激活码数量
	DeletedCnt   int         `json:"deleted_cnt" structs:"deleted_cnt"`   // 删除的激活码数量
	LockedCnt    int         `json:"locked_cnt" structs:"locked_cnt"`     // 锁定的激活码数量
	CreatedAt    time.Time   `json:"created_at" structs:"created_at"`     // 用户创建时间
	Roles        Roles       `json:"roles" structs:"roles"`               // 用户角色
	Introduction string      `json:"introduction" structs:"introduction"` // 介绍
	Avatar       string      `json:"avatar" structs:"avatar"`             // 头像
	MaxCnt       int         `json:"max_cnt" structs:"max_cnt"`           // 最大激活码数量
	Permissions  Permissions `json:"permissions" structs:"permissions"`   // 用户权限
	Apps         Apps        `json:"apps" structs:"apps"`                 // 用户有权限的应用
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) ToView() UserView {
	createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
	if u.Roles == nil {
		u.Roles = Roles{}
	}
	if u.Permissions == nil {
		u.Permissions = Permissions{}
	}
	return UserView{
		ID:           u.ID,
		Username:     u.Username,
		Status:       int(u.Status),
		Ancestry:     u.Ancestry,
		TotalCnt:     u.TotalCnt,
		UsedCnt:      u.UsedCnt,
		NousedCnt:    u.NousedCnt,
		DeletedCnt:   u.DeletedCnt,
		LockedCnt:    u.LockedCnt,
		CreatedAt:    createdAt,
		MaxCnt:       u.MaxCnt,
		Roles:        u.Roles,
		Introduction: u.Introduction,
		Avatar:       u.Avatar,
		Permissions:  u.Permissions,
		Apps:         u.Apps,
	}
}

func (u *User) ToDBStruct() DBStruct {
	return DBStruct{
		ID:           u.ID,
		Username:     u.Username,
		Password:     u.Password,
		Status:       u.Status,
		Ancestry:     u.Ancestry,
		TotalCnt:     u.TotalCnt,
		UsedCnt:      u.UsedCnt,
		NousedCnt:    u.NousedCnt,
		DeletedCnt:   u.DeletedCnt,
		LockedCnt:    u.LockedCnt,
		CreatedAt:    u.CreatedAt,
		MaxCnt:       u.MaxCnt,
		Roles:        u.Roles.ToJsonArray(),
		Introduction: u.Introduction,
		Avatar:       u.Avatar,
		Permissions:  u.Permissions.ToJsonArray(),
		Apps:         u.Apps.ToJsonArray(),
	}
}

func (u *User) IsRoot() bool {
	return u.HasRole(RoleRoot)
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) HasPermission(permission string) bool {
	// 检查状态再检查权限
	if u.Status != StatusNormal {
		return false
	}

	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (u *User) HasApp(appId string) bool {
	for _, id := range u.Apps {
		if id == appId {
			return true
		}
	}
	return false
}

func BatchToView(users []User) []UserView {
	var userViews []UserView
	for _, user := range users {
		userViews = append(userViews, user.ToView())
	}
	return userViews
}

// UserView 仅仅是为了前端展示的 User 结构
type UserView struct {
	ID           string      `json:"id"`           // 用户唯一标识符(UUID)
	Username     string      `json:"username"`     // 用户名
	Status       int         `json:"status"`       // 状态: -1-停封, 1-正常
	Ancestry     string      `json:"ancestry"`     // 祖先用户ID（如果有的话）
	TotalCnt     int         `json:"total_cnt"`    // 总激活码数量
	UsedCnt      int         `json:"used_cnt"`     // 已使用激活码数量
	NousedCnt    int         `json:"noused_cnt"`   // 未使用激活码数量
	DeletedCnt   int         `json:"deleted_cnt"`  // 删除的激活码数量
	LockedCnt    int         `json:"locked_cnt"`   // 锁定的激活码数量
	CreatedAt    string      `json:"created_at"`   // 用户创建时间
	MaxCnt       int         `json:"max_cnt"`      // 最大激活码数量
	Roles        Roles       `json:"roles"`        // 用户角色
	Introduction string      `json:"introduction"` // 介绍
	Avatar       string      `json:"avatar"`       // 头像
	Permissions  Permissions `json:"permissions"`  // 用户权限
	Apps         Apps        `json:"apps"`         // 用户有权限的应用
}
