package app

import (
	"bufio"
	"dungeon/internal/config"
	"dungeon/internal/domain/models"
	"dungeon/internal/parser"
	"dungeon/internal/service"
	"fmt"
	"os"
)

type App struct {
	config    *config.Config
	parser    *parser.EventParser
	processor *service.Processor
}

func New(configPath string) *App {
	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not load config: %v", err))
	}

	dungeon, err := models.NewDungeonFromConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("could not create dungeon: %v", err))
	}

	return &App{
		config:    cfg,
		parser:    parser.New(),
		processor: service.NewProcessor(dungeon),
	}
}

func (a *App) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if line == "" {
			continue
		}

		event, err := a.parser.Parse(line)
		if err != nil {
			return fmt.Errorf("error parsing line %d: %w", lineNum, err)
		}

		a.processor.ProcessEvent(event)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading: %w", err)
	}

	a.processor.PrintOutput()
	a.processor.PrintFinalReport()

	return nil
}
