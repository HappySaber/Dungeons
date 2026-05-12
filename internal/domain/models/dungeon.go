package models

import (
	"dungeon/internal/config"
	"fmt"
	"time"
)

type Dungeon struct {
	Floors           int
	MonstersPerFloor int
	OpenAt           time.Time
	CloseAt          time.Time
}

func NewDungeonFromConfig(cfg *config.Config) (*Dungeon, error) {
	openTime, err := ParseOpenTime(cfg.OpenAt)
	if err != nil {
		return nil, fmt.Errorf("could not parse open time: %w", err)
	}

	closeTime := CalculateCloseTime(cfg.Duration, openTime)

	return &Dungeon{
		Floors:           cfg.Floors,
		MonstersPerFloor: cfg.Monsters,
		OpenAt:           openTime,
		CloseAt:          closeTime,
	}, nil
}

func (d *Dungeon) IsOpen(t time.Time) bool {
	return !t.Before(d.OpenAt) && t.Before(d.CloseAt)
}

func (d *Dungeon) IsBossFloor(floor int) bool {
	return floor == d.Floors
}

func (d *Dungeon) IsValidFloor(floor int) bool {
	return floor >= 1 && floor <= d.Floors
}

func (d *Dungeon) TotalMonsters() int {
	return d.MonstersPerFloor * (d.Floors - 1)
}

func (d *Dungeon) MonstersOnFloor(floor int) int {
	if d.IsBossFloor(floor) {
		return 0
	}
	return d.MonstersPerFloor
}

func ParseOpenTime(openAt string) (time.Time, error) {
	return time.Parse("15:04:05", openAt)
}

func CalculateCloseTime(duration int, openTime time.Time) time.Time {
	return openTime.Add(time.Duration(duration) * time.Hour)
}
