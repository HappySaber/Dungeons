package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Floors <= 0 {
		return fmt.Errorf("floors must be > 0, got: %d", c.Floors)
	}
	if c.Monsters < 0 {
		return fmt.Errorf("monsters must be >= 0, got: %d", c.Monsters)
	}
	if c.Duration <= 0 {
		return fmt.Errorf("duration must be > 0, got: %d", c.Duration)
	}

	if _, err := time.Parse("15:04:05", c.OpenAt); err != nil {
		return fmt.Errorf("invalid format for OpenAt: %w", err)
	}

	return nil
}

func (c *Config) ParseOpenTime() (time.Time, error) {
	return time.Parse("15:04:05", c.OpenAt)
}

func (c *Config) CalculateCloseTime(openTime time.Time) time.Time {
	return openTime.Add(time.Duration(c.Duration) * time.Hour)
}
