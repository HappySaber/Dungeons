package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Floors   int `json:"Floors"`
	Monsters int `json:"Monsters"`

	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file doens't exist: " + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (c *Config) Validate() error {
	if c.Floors <= 0 {
		return fmt.Errorf("floors must be positive")
	}

	if c.Monsters <= 0 {
		return fmt.Errorf("monsters must be positive")
	}

	if _, err := time.Parse("15:04:05", c.OpenAt); err != nil {
		return fmt.Errorf("invalid OpenAt format")
	}

	if c.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	return nil
}
