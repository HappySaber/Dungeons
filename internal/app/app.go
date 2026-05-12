package app

import (
	"dungeon/internal/config"
)

type App struct {
	config *config.Config
}

func New(configPath string) *App {
	var cfg *config.Config
	cfg = config.MustLoad(configPath)

	return &App{
		config: cfg,
	}

}
