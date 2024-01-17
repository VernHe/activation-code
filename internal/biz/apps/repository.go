package apps

type Repository interface {
	QueryAppList(args QueryAppListArgs) (QueryAppListResult, error)
	CreateApp(app App) error
	UpdateApp(app App) error
	DeleteApp(id string) error
	QueryAppOptions() ([]AppOption, error)
	GetAppByIDs(ids []string) ([]App, error)
}
