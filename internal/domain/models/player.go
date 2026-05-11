package models

import "time"

type Player struct {
	ID int

	Registered bool
	InDungeon  bool
	Dead       bool
	Disqual    bool

	HP int

	CurrentFloor int

	MonstersKilledOnFloor int

	BossKilled bool

	EnterTime time.Time
	ExitTime  time.Time

	FloorEnterTime time.Time

	FloorDurations []time.Duration
	BossDuration   time.Duration
}
