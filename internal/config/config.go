package config

import "time"

type Config struct {
	Floors           int `json:"floors"`
	MonstersPerFloor int `json:"monsters_per_floor"`

	OpenAt   time.Time     `json:"open_at"`
	Duration time.Duration `json:"duration"`
}
