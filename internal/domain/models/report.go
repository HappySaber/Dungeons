package models

import "time"

type Report struct {
	State string

	PlayerID int

	TotalTime    time.Duration
	AvgFloorTime time.Duration
	BossTime     time.Duration

	HP int
}
