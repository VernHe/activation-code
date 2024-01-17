package card

type Repository interface {
	GetCardByID(id string) (Card, error)
	GetCardByValue(value string) (Card, error)
	GetCardsByUserId(userId string) ([]Card, error)
	GetCards(args GetCardsArgs) (GetCardsResult, error)
	DeleteCardByValue(value string, userId string) error
	DeleteCardsByValues(values []string, userId string) error
	CreateCard(card Card) (Card, error)
	CreateCards(cards []Card) ([]Card, error)
	UpdateCard(card Card) error
	DeleteCard(card Card) error
	GetCardCountByUserIds(userIds []string) (map[string]CardCountByUser, error)
	GetCardTotalCountByUserId(userId string) (int64, error)
	BatchUpdateStatus(args BatchUpdateStatusArgs) error
	GetCardCountByUserIdAndStatus(userId string) (map[int]int, error)
}
