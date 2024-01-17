package apps

import "time"

type App struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CardLength int       `json:"card_length"`
	CardPrefix string    `json:"card_prefix"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *App) TableName() string {
	return "app"
}

func (a *App) ToOption() AppOption {
	return AppOption{
		ID:   a.ID,
		Name: a.Name,
	}
}

func BatchToOption(apps []App) []AppOption {
	var options []AppOption
	for _, app := range apps {
		options = append(options, app.ToOption())
	}
	return options
}
