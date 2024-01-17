package user

import (
	"errors"

	"configuration-management/global"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

/*
表结构如下：
CREATE TABLE user (
    id INT PRIMARY KEY AUTO_INCREMENT, -- 用户唯一标识符
    username VARCHAR(255) NOT NULL UNIQUE, -- 用户名
    password VARCHAR(255) NOT NULL, -- 密码
    status INT DEFAULT 1, -- 状态: 0-停封, 1-正常
    ancestry VARCHAR(255), -- 祖先用户ID（如果有的话）
    ancestry_depth INT DEFAULT 0, -- 用户在家谱中的深度
    total_cnt INT DEFAULT 0, -- 总激活码数量
    used_cnt INT DEFAULT 0, -- 已使用激活码数量
    noused_cnt INT DEFAULT 0, -- 未使用激活码数量
    deleted_cnt INT DEFAULT 0, -- 删除的激活码数量
    locked_cnt INT DEFAULT 0, -- 锁定的激活码数量
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 用户创建时间
    acentry INT -- 祖先用户ID（可能是一个整数）
	permissions json -- 用户权限
);
*/

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(id string) (User, error) {
	var dbStruct DBStruct
	if err := r.db.Table((&DBStruct{}).TableName()).Where("id = ?", id).First(&dbStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("查询时记录不存在", "error:", err, "user_id", id)
			return User{}, err
		}
	}

	user, err := dbStruct.ToModel()
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"id":     id,
			"error:": err,
		}).Error("查询时记录转换失败")
		return User{}, errcode.ServerError
	}
	return user, nil
}

func (r *repository) QueryUserList(args QueryUserListArgs) (QueryUserListResult, error) {
	db := r.db.Table((&DBStruct{}).TableName())
	if args.Username != "" {
		db.Where("username like ?", "%"+args.Username+"%")
	}
	if args.Status != StatusUnknown {
		db.Where("status = ?", args.Status)
	}

	// 获取数量
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return QueryUserListResult{}, err
	}

	if args.Page != 0 {
		db.Offset((args.Page - 1) * args.Limit)
	}
	if args.Limit != 0 {
		db.Limit(args.Limit)
	}

	// order by created_at desc
	db.Order("created_at desc")
	var dbStructs = make([]DBStruct, 0)
	if err := db.Find(&dbStructs).Error; err != nil {
		global.Logger.WithFields(logger.Fields{
			"error:": err,
		}).Error("查询时记录失败")
		return QueryUserListResult{}, err
	}
	users, err := BatchToModel(dbStructs)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"error:": err,
		}).Error("查询时记录转换失败")
		return QueryUserListResult{}, err
	}
	return QueryUserListResult{
		Total: int(total),
		List:  users,
	}, nil
}

func (r *repository) GetUserByUsername(username string) (User, error) {
	var dbStruct DBStruct
	if err := r.db.Table(dbStruct.TableName()).Where("username = ?", username).First(&dbStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"username": username,
				"error:":   err,
			}).Error("查询时记录不存在")
			return User{}, err
		}
		return User{}, err
	}

	user, err := dbStruct.ToModel()
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"username": username,
			"error:":   err,
		}).Error("查询时记录转换失败")
		return User{}, errcode.ServerError
	}
	return user, nil
}

func (r *repository) CreateUser(user User) error {
	dbStruct := user.ToDBStruct()

	if err := r.db.Table((&DBStruct{}).TableName()).Debug().Create(&dbStruct).Error; err != nil {
		global.Logger.WithFields(logger.Fields{
			"user_id": user.ID,
		}).Error("创建时记录失败", err)
		return err
	}
	return nil
}

func (r *repository) UpdateUser(user User) error {
	userMap := structs.Map(user)
	// 对 []string 类型对 roles 和 permissions 字段进行特殊处理
	if roles, ok := userMap["roles"]; ok {
		userMap["roles"] = roles.(Roles).ToJsonArray()
	}
	if permissions, ok := userMap["permissions"]; ok {
		userMap["permissions"] = permissions.(Permissions).ToJsonArray()
	}
	if apps, ok := userMap["apps"]; ok {
		userMap["apps"] = apps.(Apps).ToJsonArray()
	}
	if err := r.db.Model(&user).Where("id = ?", user.ID).Updates(userMap).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id": user.ID,
			}).Error("更新时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}

func (r *repository) DeleteUser(user User) error {
	dbStruct := user.ToDBStruct()
	// 根据 ID 进行删除
	if err := r.db.Table((&dbStruct).TableName()).Where("id = ?", user.ID).Delete(&dbStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id": user.ID,
			}).Error("删除时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}
