package card

import (
	"errors"

	"configuration-management/global"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

/*
CREATE TABLE card (
    id VARCHAR(36) PRIMARY KEY, -- 激活码唯一标识符(UUID)，使用字符串表示
    status ENUM('未使用', '已使用', '已锁定', '已删除') NOT NULL, -- 状态: 0-未使用, 1-已使用, 2-已锁定, 3-已删除
    user_id VARCHAR(36) NOT NULL, -- 用户ID，关联到用户表中的id字段
    days INT NOT NULL, -- 有效天数
    expired_at TIMESTAMP NOT NULL, -- 过期时间
    value VARCHAR(255) NOT NULL, -- 激活码值
    used BOOLEAN NOT NULL, -- 是否使用过，使用布尔类型表示
    used_at TIMESTAMP, -- 使用时间
    deleted_at TIMESTAMP, -- 删除时间
    created_at TIMESTAMP NOT NULL -- 生成时间
);
*/

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCardByID(id string) (Card, error) {
	var card Card
	if err := r.db.Where("id = ?", id).First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("查询时记录不存在", "error:", err, "card_id", id)
			return Card{}, err
		}
	}
	return card, nil
}

func (r *repository) GetCardByValue(value string) (Card, error) {
	var card Card
	if err := r.db.Where("value = ?", value).First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"value":  value,
				"error:": err,
			}).Error("查询时记录不存在")
			return Card{}, errcode.NotFound
		}
		return Card{}, err
	}
	return card, nil
}

func (r *repository) GetCardsByUserId(userId string) ([]Card, error) {
	var cards []Card
	if err := r.db.Where("user_id = ?", userId).Find(&cards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id": userId,
				"error:":  err,
			}).Error("查询时记录不存在")
			return []Card{}, errcode.NotFound
		}
		return []Card{}, err
	}
	return cards, nil
}

func (r *repository) GetCards(args GetCardsArgs) (GetCardsResult, error) {
	var cards []Card
	// 动态的构建查询条件
	db := r.db.Table((&Card{}).TableName())
	if args.UserId != "" {
		db = db.Where("user_id = ?", args.UserId)
	}
	if len(args.Values) > 0 {
		db = db.Where("value IN (?)", args.Values)
	}
	if len(args.Status) > 0 {
		db = db.Where("status IN (?)", args.Status)
	}
	if args.Remark != "" {
		// 模糊查询
		db = db.Where("remark LIKE ?", "%"+args.Remark+"%")
	}
	if len(args.AppIDs) > 0 {
		db = db.Where("app_id IN (?)", args.AppIDs)
	}
	if args.TimeType != "" {
		db = db.Where("time_type = ?", args.TimeType)
	}
	if args.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+args.UserName+"%")
	}
	if args.SEID != "" {
		db = db.Where("seid = ?", args.SEID)
	}
	if !args.CreatedAtDateRange.StartTime.IsZero() {
		db = db.Where("created_at >= ?", args.CreatedAtDateRange.StartTime.Format("2006-01-02 15:04:05"))
	}
	if !args.CreatedAtDateRange.EndTime.IsZero() {
		db = db.Where("created_at <= ?", args.CreatedAtDateRange.EndTime.Format("2006-01-02 15:04:05"))
	}
	if !args.UsedAtDateRange.StartTime.IsZero() {
		db = db.Where("used_at >= ?", args.UsedAtDateRange.StartTime.Format("2006-01-02 15:04:05"))
	}
	if !args.UsedAtDateRange.EndTime.IsZero() {
		db = db.Where("used_at <= ?", args.UsedAtDateRange.EndTime.Format("2006-01-02 15:04:05"))
	}

	// 打印 SQL 语句
	db = db.Debug()

	if args.Page == 0 || args.Limit == 0 {
		args.Page = 1
		args.Limit = 10
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Logger.WithFields(logger.Fields{
			"args":   args,
			"error:": err,
		}).Error("查询时记录总数失败")
		return GetCardsResult{}, err
	}

	if args.NeedPagination {
		db = db.Offset((args.Page - 1) * args.Limit).Limit(args.Limit)
	}

	// 查询列表
	if err := db.Order("created_at DESC").Find(&cards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args":   args,
				"error:": err,
			}).Error("查询时记录不存在")
			return GetCardsResult{}, errcode.NotFound
		}
		return GetCardsResult{}, err
	}

	return GetCardsResult{
		Total: int(total),
		List:  cards,
	}, nil
}

func (r *repository) DeleteCardByValue(value string, userId string) error {
	db := r.db.Table((&Card{}).TableName()).Where("value = ?", value)
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	if err := db.Unscoped().Delete(&Card{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"value":  value,
				"error:": err,
			}).Error("删除时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}

func (r *repository) CreateCard(card Card) (Card, error) {
	// 检查重复记录
	if _, err := r.GetCardByValue(card.Value); err != nil {
		if !errors.Is(err, errcode.NotFound) {
			global.Logger.WithFields(logger.Fields{
				"card": card,
			}).Error("创建时查询记录失败", err)
			return Card{}, errcode.DuplicateKey
		}
	}

	// 创建记录
	if err := r.db.Create(&card).Error; err != nil {
		global.Logger.WithFields(logger.Fields{
			"card": card,
		}).Error("创建时记录失败", err)
		return Card{}, err
	}

	return card, nil
}

// CreateCards 根据数量批量创建 Card 记录
func (r *repository) CreateCards(cards []Card) ([]Card, error) {
	// gorm 的批量创建
	if err := r.db.Create(cards).Error; err != nil {
		global.Logger.WithFields(logger.Fields{
			"cards": cards,
		}).Error("创建时记录失败", err)
		return []Card{}, err
	}

	// 返回创建的记录
	return cards, nil
}

func (r *repository) UpdateCard(card Card) error {
	cardMap := structs.Map(card)
	if card.ExpiredAt != nil && !card.ExpiredAt.IsZero() {
		cardMap["expired_at"] = card.ExpiredAt.Format("2006-01-02 15:04:05")
	}
	if err := r.db.Model(&card).Where("id = ?", card.ID).Updates(cardMap).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"card": card,
			}).Error("更新时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}

func (r *repository) DeleteCard(card Card) error {
	// 根据 ID 进行删除
	if err := r.db.Table((&Card{}).TableName()).Delete(&Card{}, card.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *repository) DeleteCardsByValues(values []string, userId string) error {
	db := r.db.Table((&Card{}).TableName()).Where("value IN (?)", values)
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	// 根据 ID 进行删除，硬删除
	if err := db.Unscoped().Delete(&Card{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"values": values,
			}).Error("删除时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}

type CardCountByUser struct {
	UserId  string `json:"user_id"`
	Total   int64  `json:"total"`
	Unused  int64  `json:"unused"`
	Used    int64  `json:"used"`
	Locked  int64  `json:"locked"`
	Deleted int64  `json:"deleted"`
}

// GetCardCountByUserIds 根据 user_ids 查询每个人的激活码总数以及各个状态的激活码数量
func (r *repository) GetCardCountByUserIds(userIds []string) (map[string]CardCountByUser, error) {
	// 查询总数
	var cardCountByUsers []CardCountByUser
	if err := r.db.Table((&Card{}).TableName()).Select("user_id, COUNT(*) AS total, SUM(IF(status = 1, 1, 0)) AS unused, SUM(IF(status = 2, 1, 0)) AS used, SUM(IF(status = 3, 1, 0)) AS locked, SUM(IF(status = 4, 1, 0)) AS deleted").Where("user_id IN (?)", userIds).Group("user_id").Find(&cardCountByUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_ids": userIds,
			}).Error("查询时记录不存在")
			return nil, errcode.NotFound
		}
		return nil, err
	}

	// 将查询结果转换为 map 方便查询
	cardCountByUserMap := make(map[string]CardCountByUser)
	for _, cardCountByUser := range cardCountByUsers {
		cardCountByUserMap[cardCountByUser.UserId] = cardCountByUser
	}

	return cardCountByUserMap, nil
}

func (r *repository) GetCardTotalCountByUserId(userId string) (int64, error) {
	var totalCount int64
	if err := r.db.Table((&Card{}).TableName()).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id": userId,
			}).Error("查询时记录不存在")
			return 0, errcode.NotFound
		}
		return 0, err
	}
	return totalCount, nil
}

// BatchUpdateStatus 批量跟新状态
func (r *repository) BatchUpdateStatus(args BatchUpdateStatusArgs) error {
	db := r.db.Table((&Card{}).TableName()).Where("value IN (?)", args.Values)

	// 如果是更新为未使用，则重置使用时间、过期时间、IMEI、SEID
	// 更新时需要更新时间
	// 由于不会更新零值字段，所以需要用sql语句更新
	var sql string
	switch args.Status {
	case StatusUnused:
		// 如果是更新为未使用，则重置使用时间、过期时间、SEID
		if args.UserId != "" {
			sql = `UPDATE card SET status = ?, used_at = NULL, expired_at = NULL, seid = NULL WHERE user_id = ? AND value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.UserId, args.Values, StatusUnused)
		} else {
			sql := `UPDATE card SET status = ?, used_at = NULL, expired_at = NULL, seid = NULL WHERE value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.Values, StatusUnused)
		}
	case StatusUsed:
		// 如果是更新为已使用，则更新使用时间、过期时间、SEID
		if args.UserId != "" {
			sql = `UPDATE card SET status = ?, used_at = NOW(), expired_at = IF(days > 0, DATE_ADD(NOW(), INTERVAL days DAY), NOW()) + INTERVAL IF(minutes > 0, minutes, 0) MINUTE, seid = ? WHERE user_id = ? AND value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.SEID, args.UserId, args.Values, StatusUsed)
		} else {
			sql = `UPDATE card SET status = ?, used_at = NOW(), expired_at = IF(days > 0, DATE_ADD(NOW(), INTERVAL days DAY), NOW()) + INTERVAL IF(minutes > 0, minutes, 0) MINUTE, seid = ? WHERE value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.SEID, args.Values, StatusUsed)
		}
		// 如果是更新为已锁定，则更新锁定时间
	case StatusLocked:
		if args.UserId != "" {
			sql = `UPDATE card SET status = ?, locked_at = NOW() WHERE user_id = ? AND value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.UserId, args.Values, StatusLocked)
		} else {
			sql = `UPDATE card SET status = ?, locked_at = NOW() WHERE value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.Values, StatusLocked)
		}
		// 如果是更新为已删除，则更新删除时间
	case StatusDeleted:
		if args.UserId != "" {
			sql = `UPDATE card SET status = ?, deleted_at = NOW() WHERE user_id = ? AND value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.UserId, args.Values, StatusDeleted)
		} else {
			sql = `UPDATE card SET status = ?, deleted_at = NOW() WHERE value IN (?) AND status != ?`
			db = db.Exec(sql, args.Status, args.Values, StatusDeleted)
		}
	}

	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"args": args,
			}).Error("更新时记录不存在")
			return errcode.NotFound
		}
		return err
	}
	return nil
}

// GetCardCountByUserIdAndStatus 统计某个用户的已使用、未使用、已锁定、已删除的激活码数量
func (r *repository) GetCardCountByUserIdAndStatus(userId string) (map[int]int, error) {
	// SELECT
	//    status,
	//    COUNT(*) as status_count
	//FROM
	//    card
	//GROUP BY
	//    status;
	var cardCountByStatus []struct {
		Status int `json:"status"`
		Count  int `json:"count"`
	}
	db := r.db.Table((&Card{}).TableName()).Select("status, COUNT(*) AS count")
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	if err := db.Group("status").Find(&cardCountByStatus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id": userId,
			}).Error("查询时记录不存在")
			return nil, errcode.NotFound
		}
		return nil, err
	}

	cardCountByStatusMap := make(map[int]int)
	for _, cardCountByStatus := range cardCountByStatus {
		cardCountByStatusMap[cardCountByStatus.Status] = cardCountByStatus.Count
	}
	return cardCountByStatusMap, nil
}
