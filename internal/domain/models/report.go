package models

import (
	"time"
)

type Report struct {
	State        FinalState
	PlayerID     int
	TotalTime    time.Duration
	AvgFloorTime time.Duration
	BossTime     time.Duration
	FinalHealth  int
}
