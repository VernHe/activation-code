package card

import (
	"errors"
	"sync"
	"time"

	"configuration-management/global"
	"configuration-management/internal/biz/apps"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"
	"configuration-management/utils"

	"github.com/patrickmn/go-cache"
)

var (
	initializing     sync.Once
	activateCache    cache.Cache       // 激活激活码缓存
	activateCacheTTL = time.Minute * 5 // 激活激活码缓存过期时间(5分钟内不能重复激活)

	checkCardStatusCache    cache.Cache       // 检查激活码状态缓存
	checkCardStatusCacheTTL = time.Minute * 5 // 检查激活码状态缓存过期时间(5分钟内更新一次)
)

type service struct {
	repo    Repository
	appRepo apps.Repository
}

func NewService() Service {
	initializing.Do(func() {
		activateCache = *cache.New(activateCacheTTL, activateCacheTTL)
		checkCardStatusCache = *cache.New(checkCardStatusCacheTTL, checkCardStatusCacheTTL)
	})
	return &service{
		repo:    NewRepository(global.DBEngine),
		appRepo: apps.NewRepository(),
	}
}

func (s *service) GetCardByID(id string) (Card, error) {
	return s.repo.GetCardByID(id)
}

func (s *service) GetCardByValue(value string) (Card, error) {
	return s.repo.GetCardByValue(value)
}

func (s *service) GetCardsByUserId(userId string) ([]Card, error) {
	return s.repo.GetCardsByUserId(userId)
}

func (s *service) GetCards(args GetCardsArgs) (GetCardsResult, error) {
	return s.repo.GetCards(args)
}

func (s *service) DeleteCardByValue(value string, userId string) error {
	return s.repo.DeleteCardByValue(value, userId)
}

func (s *service) CreateCard(args CreateCardArgs) (Card, error) {
	now := time.Now()
	var card = Card{
		ID:        utils.GenerateUUID(),
		Status:    StatusUnused,
		UserID:    args.UserID,
		Days:      args.Days,
		Value:     utils.GenerateActivationKey(),
		CreatedAt: &now,
	}

	newCard, err := s.repo.CreateCard(card)
	if err != nil {
		if errors.Is(err, errcode.DuplicateKey) {
			// 重复则重新生成Value
			card.Value = utils.GenerateActivationKey()
			return s.repo.CreateCard(card)
		}
		return Card{}, err
	}

	return newCard, nil
}

func (s *service) CreateCards(args CreateCardsArgs) ([]Card, error) {
	// 校验AppID
	result, err := s.appRepo.QueryAppList(apps.QueryAppListArgs{
		ID: args.AppID,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("查询应用列表失败", err)
		return []Card{}, err
	}
	if result.Total == 0 {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("应用不存在")
		return []Card{}, errcode.NotFound.WithDetails("应用不存在")
	}

	// 获取 App 信息
	appList, err := s.appRepo.QueryAppList(apps.QueryAppListArgs{
		ID: args.AppID,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("查询应用列表失败", err)
		return []Card{}, err
	}
	if appList.Total == 0 {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("应用不存在")
		return []Card{}, errcode.NotFound.WithDetails("应用不存在")
	}
	app := appList.List[0]

	// 检查用户最大创建数量
	totalCnt, err := s.repo.GetCardTotalCountByUserId(args.UserID)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("查询用户激活码数量失败", err)
		return []Card{}, err
	}
	if totalCnt+int64(args.Count) > int64(args.MaxCnt) {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("用户激活码数量超限")
		return []Card{}, errcode.NoPermission.WithDetails("用户激活码数量超限")
	}

	now := time.Now()
	var cards []Card
	for i := 0; i < args.Count; i++ {
		cards = append(cards, Card{
			ID:        utils.GenerateUUID(),
			AppID:     args.AppID,
			Status:    StatusUnused,
			UserID:    args.UserID,
			UserName:  args.UserName,
			Minutes:   args.Minutes,
			TimeType:  args.TimeType,
			Days:      args.Days,
			Value:     utils.GenerateActivationKeyByApp(app.CardPrefix, app.CardLength),
			Remark:    args.Remark,
			CreatedAt: &now,
		})
	}

	newCards, err := s.repo.CreateCards(cards)
	if err != nil {
		return []Card{}, err
	}

	return newCards, nil
}

func (s *service) UpdateCard(args UpdateCardArgs) error {
	// 非 root 用户
	if args.UserId != "" {
		// 检查用户是否有权限
		if args.CurrentCard.UserID != args.UserId {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Error("用户无权限")
			return errcode.NoPermission.WithDetails("用户无权限")
		}
	}

	// 更新字段
	args.CurrentCard.Status = args.Status
	args.CurrentCard.SEID = args.SEID
	args.CurrentCard.Minutes = args.Minutes
	args.CurrentCard.Days = args.Days
	args.CurrentCard.TimeType = args.TimeType
	args.CurrentCard.Remark = args.Remark

	return s.repo.UpdateCard(args.CurrentCard)
}

func (s *service) DeleteCard(card Card) error {
	// 将状态更新为已删除
	c, err := s.repo.GetCardByID(card.ID)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"card": card,
		}).Error("更新时查询记录失败", err)
		return err
	}

	c.Status = StatusDeleted

	return s.repo.UpdateCard(c)
}

func (s *service) DeleteCardsByValues(values []string, userId string) error {
	return s.repo.DeleteCardsByValues(values, userId)
}

func (s *service) BatchUpdateStatus(args BatchUpdateStatusArgs) error {
	return s.repo.BatchUpdateStatus(args)
}

func (s *service) GetCardCountByUserIdAndStatus(userId string) (map[int]int, error) {
	return s.repo.GetCardCountByUserIdAndStatus(userId)
}

// ActivateCard 激活激活码
func (s *service) ActivateCard(args ActivateCardArgs) (Card, error) {
	// 检查缓存，防止重复激活
	if c, ok := activateCache.Get(args.Value); ok {
		// 直接返回结果
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("激活码被重复激活")
		return c.(Card), nil
	}

	// 校验激活码是否存在
	card, err := s.repo.GetCardByValue(args.Value)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("查询激活码失败", err)
		return Card{}, err
	}
	if card.ID == "" {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("激活码不存在")
		return Card{}, errcode.NotFound.WithDetails("激活码不存在")
	}

	// 激活条件: 未使用且未过期且未绑定设备

	// 重复激活
	if card.Status == StatusUsed {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("激活码被重复激活")
		return card, nil
	}

	// 校验激活码状态
	if card.Status != StatusUnused {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("激活码状态不正确")
		return card, errors.New("激活码状态不正确")
	}

	// 校验是否有设备绑定
	if args.SEID == "" {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("激活码未绑定设备")
		return card, errcode.InvalidParams.WithDetails("激活码缺少绑定设备")
	}

	// 更新激活码状态
	card.Status = StatusUsed
	card.SEID = args.SEID
	card.Used = true
	now := time.Now()
	card.UsedAt = &now
	// Days + Minutes
	expiredAt := now.AddDate(0, 0, card.Days).Add(time.Minute * time.Duration(card.Minutes))
	card.ExpiredAt = &expiredAt

	err = s.repo.UpdateCard(card)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("更新激活码失败", err)
		return Card{}, err
	}

	// 设置缓存
	activateCache.SetDefault(args.Value, card)

	return card, nil
}

// CheckCardStatus 检查激活码目前是否处于能使用的状态
func (s *service) CheckCardStatus(args CheckCardStatusArgs) (bool, error) {
	// 检查缓存，提前返回结果
	result, ok := checkCardStatusCache.Get(args.Value)
	if ok {
		println("命中缓存")
		return result.(bool), nil
	}

	// 校验激活码是否存在
	card, err := s.repo.GetCardByValue(args.Value)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"value": args.Value,
		}).Error("查询激活码失败", err)
		if errors.Is(err, errcode.NotFound) {
			checkCardStatusCache.SetDefault(args.Value, false)
			return false, nil
		}
		return false, err
	}
	if card.ID == "" {
		global.Logger.WithFields(logger.Fields{
			"value": args.Value,
		}).Error("激活码不存在")
		checkCardStatusCache.SetDefault(args.Value, false)
		return false, nil
	}

	// 校验激活码状态
	if card.Status != StatusUsed {
		global.Logger.WithFields(logger.Fields{
			"card": card,
		}).Error("激活码未激活")
		return false, nil
	}

	// 校验激活码是否过期
	if card.ExpiredAt.IsZero() || card.ExpiredAt.Before(time.Now()) {
		global.Logger.WithFields(logger.Fields{
			"card": card,
		}).Error("激活码已过期")
		return false, nil
	}

	// 检查设备是否匹配
	if card.SEID != args.SEID {
		global.Logger.WithFields(logger.Fields{
			"card": card,
			"args": args,
		}).Error("设备不匹配")
		return false, nil
	}

	// 设置缓存
	checkCardStatusCache.SetDefault(args.Value, true)

	// 状态为已激活且未过期且设备匹配
	return true, nil
}

// SetCardExpiredAt 增加激活码的时间
func (s *service) SetCardExpiredAt(args SetCardExpiredAtArgs) error {
	// 校验激活码是否存在
	card, err := s.repo.GetCardByValue(args.Value)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"value": args.Value,
		}).Error("查询激活码失败", err)
		return err
	}
	if card.ID == "" {
		global.Logger.WithFields(logger.Fields{
			"value": args.Value,
		}).Error("激活码不存在")
		return errcode.NotFound.WithDetails("激活码不存在")
	}

	if args.UserId != "" {
		// 检查用户是否有权限
		if card.UserID != args.UserId {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Error("用户无权限")
			return errcode.NoPermission.WithDetails("用户无权限")
		}
	}

	// 校验激活码状态
	if card.Status != StatusUsed {
		global.Logger.WithFields(logger.Fields{
			"card": card,
		}).Error("激活码未激活")
		return errcode.NoPermission.WithDetails("激活码未激活")
	}

	// 更新激活码状态
	card.ExpiredAt = &args.NewExpiredAt

	return s.repo.UpdateCard(card)
}
