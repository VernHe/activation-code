package user

type Repository interface {
	GetUserByID(id string) (User, error)
	QueryUserList(args QueryUserListArgs) (QueryUserListResult, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(user User) error
}
