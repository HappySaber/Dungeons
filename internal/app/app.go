package app

import (
	"dungeon/internal/config"
	"dungeon/internal/parser"
)

type App struct {
	config *config.Config
	parser *parser.EventParser
}

func New(configPath string) *App {
	var cfg *config.Config
	cfg = config.MustLoad()

	return &App{
		config: cfg,
		parser: parser.New(),
	}

}
