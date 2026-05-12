package app

import (
	"dungeon/internal/config"
)

type App struct {
	config *config.Config
}

func New(configPath string) *App {

	return &App{}
}
