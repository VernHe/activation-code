package apps

type QueryAppListArgs struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type QueryAppListResult struct {
	List  []App
	Total int
}

type CreateAppArgs struct {
	Name       string `json:"name"`
	CardLength int    `json:"card_length"`
	CardPrefix string `json:"card_prefix"`
}

type UpdateAppArgs struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CardLength int    `json:"card_length"`
	CardPrefix string `json:"card_prefix"`
}

type AppOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	QueryAppList(args QueryAppListArgs) (QueryAppListResult, error)
	CreateApp(args CreateAppArgs) error
	UpdateApp(args UpdateAppArgs) error
	DeleteApp(id string) error
	QueryAppOptions() ([]AppOption, error)
	GetAppByIDs(ids []string) ([]App, error)
}
