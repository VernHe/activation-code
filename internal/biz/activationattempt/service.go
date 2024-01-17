package activationattempt

type Service interface {
	CreateActivationAttempt(attempt *ActivationAttempt) error
	GetActivationAttemptCountByCardValue(cardValue string) (int64, error)
	GetActivationAttemptCountByCardValueInHour(cardValue string) (int64, error)
}
