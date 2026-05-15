package service

import (
	"dungeon/internal/domain/models"
	"time"
)

type Reporter struct {
	dungeon *models.Dungeon
}

func NewReporter(dungeon *models.Dungeon) *Reporter {
	return &Reporter{
		dungeon: dungeon,
	}
}

// Generate creates a Report for the given Player by analyzing their state, time spent in the dungeon, monsters killed,
// and other relevant metrics to determine their final outcome (success, fail, or disqualification) and compile all this
// information into a structured report format
func (r *Reporter) Generate(player *models.Player) *models.Report {
	return &models.Report{
		State:        r.determineState(player),
		PlayerID:     player.ID,
		TotalTime:    r.calculateTotalTime(player),
		AvgFloorTime: r.calculateAvgFloorTime(player),
		BossTime:     r.calculateBossTime(player),
		FinalHealth:  player.Health,
	}
}

func (r *Reporter) determineState(player *models.Player) models.FinalState {
	if player.State == models.StateDisqualified {
		return models.FinalDisqual
	}

	if player.State == models.StateDead {
		return models.FinalFail
	}

	if !player.IsInDungeon() && player.State != models.StateCompleted {
		return models.FinalDisqual
	}

	expectedMonsters := r.dungeon.TotalMonsters()
	if player.MonstersKilled >= expectedMonsters && player.BossKilled {
		return models.FinalSuccess
	}

	return models.FinalFail
}

func (r *Reporter) calculateTotalTime(player *models.Player) time.Duration {
	if player.EnteredAt.IsZero() {
		return 0
	}

	endTime := player.LeftAt
	if endTime.IsZero() || endTime.Before(player.EnteredAt) {
		endTime = r.dungeon.CloseAt
	}

	return endTime.Sub(player.EnteredAt)
}

func (r *Reporter) calculateAvgFloorTime(player *models.Player) time.Duration {
	if len(player.FloorClearTimes) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range player.FloorClearTimes {
		total += t
	}

	return total / time.Duration(len(player.FloorClearTimes))
}

func (r *Reporter) calculateBossTime(player *models.Player) time.Duration {
	if player.BossFloorEnteredAt.IsZero() || player.BossKilledAt.IsZero() {
		return 0
	}

	return player.BossKilledAt.Sub(player.BossFloorEnteredAt)
}
