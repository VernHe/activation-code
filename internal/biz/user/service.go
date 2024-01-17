package user

type CreateUserArgs struct {
	CreatorID    string      `json:"creator"`
	Username     string      `json:"username"`
	Password     string      `json:"password"`
	MaxCnt       int         `json:"max_cnt"`
	Apps         Apps        `json:"apps"`
	Permissions  Permissions `json:"permissions"`
	Introduction string      `json:"introduction"`
}

type LoginArgs struct {
	Username    string
	PasswordMD5 string
}

type GetUserInfoArgs struct {
	UserId string `json:"user_id"`
}

type QueryUserListArgs struct {
	Username string
	Status   Status
	Page     int `json:"page"`
	Limit    int `json:"limit"`
}

type QueryUserListResult struct {
	List  []User
	Total int
}

type UpdateUserArgs struct {
	ID           string      `json:"id"`
	Status       int         `json:"status"`
	MaxCnt       int         `json:"max_cnt"`
	UpdaterID    string      `json:"updater_id"`
	Apps         Apps        `json:"apps"`
	Roles        Roles       `json:"roles"`
	Permissions  Permissions `json:"permissions"`
	Introduction string      `json:"introduction" binding:"required"`
}

type ResetPasswordArgs struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type Service interface {
	GetUserByID(id string) (User, error)
	GetUserByUsername(username string) (User, error)
	QueryUserList(args QueryUserListArgs) (QueryUserListResult, error)
	CreateUser(args CreateUserArgs) error
	UpdateUser(args UpdateUserArgs) error
	DeleteUser(user User) error
	Login(args LoginArgs) (User, error)
	GetUserInfo(args GetUserInfoArgs) (UserView, error)
	ResetPassword(args ResetPasswordArgs) error
}
