package main

import (
	"dungeon/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	err := cfg.Validate()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
