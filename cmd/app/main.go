package main

import (
	"dungeon/internal/app"
	"flag"
	"log"
	"os"
)

func main() {

	eventsPath := flag.String("events", "", "path to events file")
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	application := app.New(*configPath)

	var err error
	if *eventsPath != "" {
		err = application.ProcessFile(*eventsPath)
	} else {
		err = application.ProcessReader(os.Stdin)
	}

	if err != nil {
		log.Fatalf("Error executing: %v", err)
	}
}
