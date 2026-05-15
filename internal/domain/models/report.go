package models

import (
	timeutil "dungeon/internal/utils"
	"fmt"
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

func (r Report) String() string {
	total := r.TotalTime

	return fmt.Sprintf(
		"[%s] %d [%s, %s, %s] HP:%d",
		r.State,
		r.PlayerID,
		timeutil.FormatDuration(total),
		timeutil.FormatDuration(r.AvgFloorTime),
		timeutil.FormatDuration(r.BossTime),
		r.FinalHealth,
	)
}
