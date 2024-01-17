package apps

import (
	"time"

	"configuration-management/global"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"
	"configuration-management/utils"
)

type serviceImpl struct {
	repo Repository
}

func NewService() Service {
	return &serviceImpl{
		repo: NewRepository(),
	}
}

func (s *serviceImpl) QueryAppList(args QueryAppListArgs) (QueryAppListResult, error) {
	// 默认的 page 和 limit
	if args.Page == 0 {
		args.Page = 1
	}
	if args.Limit == 0 {
		args.Limit = 10
	}
	return s.repo.QueryAppList(args)
}

func (s *serviceImpl) CreateApp(args CreateAppArgs) error {
	result, err := s.repo.QueryAppList(QueryAppListArgs{
		Name: args.Name,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"args": args,
		}).Error("查询应用列表失败", "error", err)
		return err
	}
	if result.Total > 0 {
		return errcode.DuplicateKey
	}

	return s.repo.CreateApp(App{
		ID:         utils.GenerateUUID(),
		Name:       args.Name,
		CardLength: args.CardLength,
		CardPrefix: args.CardPrefix,
		CreatedAt:  time.Now(),
	})
}

func (s *serviceImpl) UpdateApp(args UpdateAppArgs) error {
	// 检查想要修改的 apps name 是否已经存在
	if args.Name != "" {
		result, err := s.repo.QueryAppList(QueryAppListArgs{
			Name: args.Name,
		})
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"nme": args.Name,
			}).Error("查询应用列表失败", "error", err)
			return err
		}
		// 排除自己
		for _, app := range result.List {
			// 如果发现有相同的 name, 且 id 不同, 则返回错误
			if app.ID != args.ID {
				return errcode.DuplicateKey
			}
		}
	}

	result, err := s.repo.QueryAppList(QueryAppListArgs{
		ID: args.ID,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"nme": args.Name,
		}).Error("更新的 APP 不存在", "error", err)
		return err
	}
	if result.Total < 1 {
		return errcode.NotFound
	}

	app := result.List[0]
	app.Name = args.Name
	app.CardPrefix = args.CardPrefix
	app.CardLength = args.CardLength
	return s.repo.UpdateApp(app)
}

func (s *serviceImpl) DeleteApp(id string) error {
	return s.repo.DeleteApp(id)
}

func (s *serviceImpl) QueryAppOptions() ([]AppOption, error) {
	return s.repo.QueryAppOptions()
}

func (s *serviceImpl) GetAppByIDs(ids []string) ([]App, error) {
	return s.repo.GetAppByIDs(ids)
}
