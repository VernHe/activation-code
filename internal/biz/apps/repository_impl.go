package apps

import (
	"configuration-management/global"

	"gorm.io/gorm"
)

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repositoryImpl{
		db: global.DBEngine,
	}
}

func (r *repositoryImpl) QueryAppList(args QueryAppListArgs) (QueryAppListResult, error) {
	db := r.db.Table((&App{}).TableName())
	if args.ID != "" {
		db.Where("id = ?", args.ID)
	}
	if args.Name != "" {
		// 模糊搜索
		db.Where("name like ?", "%"+args.Name+"%")
	}

	// 获取数量
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return QueryAppListResult{}, err
	}

	if args.Page != 0 {
		db.Offset((args.Page - 1) * args.Limit)
	}
	if args.Limit != 0 {
		db.Limit(args.Limit)
	}

	// order by created_at desc
	db.Order("created_at desc")
	var apps []App
	if err := db.Find(&apps).Error; err != nil {
		return QueryAppListResult{}, err
	}
	return QueryAppListResult{
		Total: int(total),
		List:  apps,
	}, nil
}

// QueryAppOptions 获取 app options, [{id:xxx, name:xxx}]
func (r *repositoryImpl) QueryAppOptions() ([]AppOption, error) {
	var apps []App
	if err := r.db.Table((&App{}).TableName()).Find(&apps).Error; err != nil {
		return nil, err
	}
	var appOptions []AppOption
	for _, app := range apps {
		appOptions = append(appOptions, AppOption{
			ID:   app.ID,
			Name: app.Name,
		})
	}
	return appOptions, nil
}

func (r *repositoryImpl) CreateApp(app App) error {
	if err := r.db.Create(&app).Error; err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateApp(app App) error {
	if err := r.db.Save(&app).Error; err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) DeleteApp(id string) error {
	if err := r.db.Delete(&App{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) GetAppByIDs(ids []string) ([]App, error) {
	var apps []App
	if err := r.db.Table((&App{}).TableName()).Where("id in ?", ids).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}
