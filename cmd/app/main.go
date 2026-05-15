package main

import (
	"dungeon/internal/app"
	"flag"
	"log"
)

func main() {

	eventsPath := flag.String("events", "data/events", "path to events file")
	configPath := flag.String("config", "config/config.json", "path to config file")
	flag.Parse()

	application := app.New(*configPath)

	var err error
	if *eventsPath != "" {
		err = application.ProcessFile(*eventsPath)
	}

	if err != nil {
		log.Fatalf("Error executing: %v", err)
	}
}
