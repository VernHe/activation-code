package activationattempt

import "configuration-management/global"

type service struct {
	Repository Repository
}

func NewService() Service {
	return &service{Repository: NewRepository(global.DBEngine)}
}

func (s *service) CreateActivationAttempt(attempt *ActivationAttempt) error {
	return s.Repository.CreateActivationAttempt(attempt)
}

func (s *service) GetActivationAttemptCountByCardValue(cardValue string) (int64, error) {
	return s.Repository.GetActivationAttemptCountByCardValue(cardValue)
}

func (s *service) GetActivationAttemptCountByCardValueInHour(cardValue string) (int64, error) {
	return s.Repository.GetActivationAttemptCountByCardValueInHour(cardValue)
}
