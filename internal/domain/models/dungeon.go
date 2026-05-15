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

// NewDungeonFromConfig creates a Dungeon instance from the provided configuration
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

// IsOpen checks if the dungeon is open at the given time
func (d *Dungeon) IsOpen(t time.Time) bool {
	return !t.Before(d.OpenAt) && t.Before(d.CloseAt)
}

// IsBossFloor checks if the given floor is the boss floor
func (d *Dungeon) IsBossFloor(floor int) bool {
	return floor == d.Floors
}

// IsValidFloor checks if the given floor number is valid within the dungeon
func (d *Dungeon) IsValidFloor(floor int) bool {
	return floor >= 1 && floor <= d.Floors
}

// TotalMonsters calculates the total number of monsters in the dungeon (excluding the boss)
func (d *Dungeon) TotalMonsters() int {
	return d.MonstersPerFloor * (d.Floors - 1)
}

// MonstersOnFloor returns the number of monsters on a given floor, accounting for the boss floor
func (d *Dungeon) MonstersOnFloor(floor int) int {
	if d.IsBossFloor(floor) {
		return 0
	}
	return d.MonstersPerFloor
}

// TimeUntilClose calculates the remaining time until the dungeon closes from the given time
func ParseOpenTime(openAt string) (time.Time, error) {
	return time.Parse("15:04:05", openAt)
}

// CalculateCloseTime computes the closing time based on the opening time and duration
func CalculateCloseTime(duration int, openTime time.Time) time.Time {
	return openTime.Add(time.Duration(duration) * time.Hour)
}
