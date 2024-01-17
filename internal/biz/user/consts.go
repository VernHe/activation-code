package user

const (
	SuperAdminUserName = "admin"

	StatusNormal  Status = 1 // 正常
	StatusUnknown Status = 0
	StatusBanned  Status = -1 // 停封

	RoleAdmin = "admin"
	RoleRoot  = "root"
)

var (
	StatusArray = []Status{
		StatusNormal,
		StatusBanned,
	}

	DefaultUserRoles = []string{
		RoleAdmin,
	}
)
