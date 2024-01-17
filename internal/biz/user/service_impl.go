package user

import (
	"errors"
	"strings"
	"time"

	"configuration-management/internal/biz/card"

	"configuration-management/global"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"
	"configuration-management/utils"

	"github.com/jinzhu/gorm"
)

type service struct {
	repo     Repository
	cardRepo card.Repository
}

func NewService() Service {
	return &service{
		repo:     NewRepository(global.DBEngine),
		cardRepo: card.NewRepository(global.DBEngine),
	}
}

func (s *service) GetUserByID(id string) (User, error) {
	return s.repo.GetUserByID(id)
}

func (s *service) GetUserByUsername(username string) (User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *service) QueryUserList(args QueryUserListArgs) (QueryUserListResult, error) {
	if args.Status != StatusUnknown {
		isValidStatus := false
		for _, status := range StatusArray {
			if args.Status == status {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			return QueryUserListResult{}, errcode.InvalidParams.WithDetails("invalid status")
		}
	}

	if args.Page <= 0 {
		args.Page = 1
	}
	if args.Limit <= 0 {
		args.Limit = 10
	}

	userList, err := s.repo.QueryUserList(args)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("[QueryUserList] 查询用户列表失败", err)
		return QueryUserListResult{}, err
	}

	// 获取ids
	userIds := make([]string, 0)
	for _, user := range userList.List {
		userIds = append(userIds, user.ID)
	}

	// 查询数量
	userCardCount, err := s.cardRepo.GetCardCountByUserIds(userIds)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("[QueryUserList] 查询用户卡片数量失败", err)
		return QueryUserListResult{}, err
	}

	for i, _ := range userList.List {
		user := &userList.List[i]
		user.TotalCnt = int(userCardCount[user.ID].Total)
		user.UsedCnt = int(userCardCount[user.ID].Used)
		user.NousedCnt = int(userCardCount[user.ID].Unused)
		user.DeletedCnt = int(userCardCount[user.ID].Deleted)
		user.LockedCnt = int(userCardCount[user.ID].Locked)
	}

	return userList, nil
}

func (s *service) CreateUser(args CreateUserArgs) error {
	// 检查被创建用户是否已经存在
	user, _ := s.repo.GetUserByUsername(args.Username)
	if user.ID != "" {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Info("[CreateUser] 用户已经存在")
		return errcode.DuplicateKey.WithDetails("用户已经存在")
	}

	creator, err := s.repo.GetUserByID(args.CreatorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[CreateUser] 创建者不存在")
			return errcode.NotFound.WithDetails("创建者不存在")
		}
		return err
	}

	// 权限检查
	if !creator.HasRole(RoleRoot) {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Info("[CreateUser] 创建者没有权限")
		return errcode.NoPermission
	}

	user = User{
		ID:           utils.GenerateUUID(),
		Username:     args.Username,
		Password:     args.Password,
		Status:       StatusNormal,
		Ancestry:     args.CreatorID,
		CreatedAt:    time.Now(),
		Apps:         args.Apps,
		Permissions:  args.Permissions,
		Roles:        DefaultUserRoles,
		MaxCnt:       args.MaxCnt,
		Introduction: args.Introduction,
	}
	return s.repo.CreateUser(user)
}

func (s *service) UpdateUser(args UpdateUserArgs) error {
	updater, err := s.repo.GetUserByID(args.UpdaterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[UpdateUser] 没有权限")
			return errcode.NotFound.WithDetails("没有权限")
		}
		return err
	}

	// 权限检查
	if !updater.IsRoot() {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Info("[UpdateUser] 没有权限")
		return errcode.NoPermission
	}

	user, err := s.repo.GetUserByID(args.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[UpdateUser] 用户不存在")
			return errcode.NotFound.WithDetails("用户不存在")
		}
		return err
	}

	if user.Username == SuperAdminUserName {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Info("[UpdateUser] 不能更新超级管理员")
		return errcode.NoPermission
	}

	user.MaxCnt = args.MaxCnt
	user.Status = NewStatus(args.Status)
	//user.Roles = args.Roles
	user.Apps = args.Apps
	user.Permissions = args.Permissions
	user.Introduction = args.Introduction

	return s.repo.UpdateUser(user)
}

func (s *service) DeleteUser(user User) error {
	return s.repo.DeleteUser(user)
}

func (s *service) Login(args LoginArgs) (User, error) {
	// 根据 id 查询
	user, err := s.repo.GetUserByUsername(args.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 登陆的用户不存在
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[Login] 用户不存在")
			return User{}, errcode.NotFound
		}
		// 查询时出现错误
		return User{}, err
	}

	// 检查用户是否被封
	if user.Status == StatusBanned {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Info("[Login] 用户已经被停封")
		return User{}, errors.New("用户已经被停封")
	}

	// 校验密码
	if !strings.EqualFold(utils.MD5(user.Password), args.PasswordMD5) {
		return User{}, errors.New("密码错误")
	}

	return user, nil
}

func (s *service) GetUserInfo(args GetUserInfoArgs) (UserView, error) {
	user, err := s.repo.GetUserByID(args.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[GetUserInfo] 查询的用户不存在")
			return UserView{}, errcode.NotFound.WithDetails(err.Error())
		}
		return UserView{}, err
	}

	return user.ToView(), nil
}

func (s *service) ResetPassword(args ResetPasswordArgs) error {
	user, err := s.repo.GetUserByID(args.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Info("[ResetPassword] 用户不存在")
			return errcode.NotFound.WithDetails("用户不存在")
		}
		return err
	}

	user.Password = args.Password
	return s.repo.UpdateUser(user)
}
