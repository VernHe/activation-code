package activationattempt

import (
	"time"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateActivationAttempt(attempt *ActivationAttempt) error {
	structMap := structs.Map(attempt)
	structMap["activation_at"] = attempt.ActivationAt.Format(time.DateTime)
	return r.db.Table((&ActivationAttempt{}).TableName()).Create(structMap).Error
}

func (r *repository) GetActivationAttemptByID(id uint) (*ActivationAttempt, error) {
	var attempt ActivationAttempt
	err := r.db.First(&attempt, id).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *repository) UpdateActivationAttempt(attempt *ActivationAttempt) error {
	structMap := structs.Map(attempt)
	structMap["activation_at"] = attempt.ActivationAt.Format(time.DateTime)
	return r.db.Model(attempt).Updates(structMap).Error
}

func (r *repository) DeleteActivationAttemptByID(id uint) error {
	return r.db.Delete(&ActivationAttempt{}, id).Error
}

// GetActivationAttemptCountByCardValue 根据 CardValue 获取累计激活次数
func (r *repository) GetActivationAttemptCountByCardValue(cardValue string) (int64, error) {
	var count int64
	err := r.db.Model(&ActivationAttempt{}).Where("card_value = ?", cardValue).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetActivationAttemptCountByCardValueInHour 获取一小时内激活的次数
func (r *repository) GetActivationAttemptCountByCardValueInHour(cardValue string) (int64, error) {
	var count int64
	err := r.db.Model(&ActivationAttempt{}).Where("card_value = ? AND activation_at > ?", cardValue, time.Now().Add(-time.Hour)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
