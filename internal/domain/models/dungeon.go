package models

import "time"

type Dungeon struct {
	Floors           int
	MonstersPerFloor int

	OpenAt   time.Time
	Duration time.Duration
}
