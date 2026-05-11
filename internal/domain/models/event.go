package models

import "time"

type Event struct {
	Time       time.Time
	ID         int
	PlayerID   int
	ExtraParam string
}
