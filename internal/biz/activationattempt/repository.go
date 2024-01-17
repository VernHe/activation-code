package activationattempt

type Repository interface {
	CreateActivationAttempt(attempt *ActivationAttempt) error
	GetActivationAttemptByID(id uint) (*ActivationAttempt, error)
	UpdateActivationAttempt(attempt *ActivationAttempt) error
	DeleteActivationAttemptByID(id uint) error
	GetActivationAttemptCountByCardValue(cardID string) (int64, error)
	GetActivationAttemptCountByCardValueInHour(cardValue string) (int64, error)
}
