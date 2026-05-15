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

// Load reads the configuration from a JSON file and validates it
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

// Validate checks if the configuration values are valid
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

// ParseOpenTime converts the OpenAt string to a time.Time object
func (c *Config) ParseOpenTime() (time.Time, error) {
	return time.Parse("15:04:05", c.OpenAt)
}

// CalculateCloseTime computes the closing time based on the opening time and duration
func (c *Config) CalculateCloseTime(openTime time.Time) time.Time {
	return openTime.Add(time.Duration(c.Duration) * time.Hour)
}
