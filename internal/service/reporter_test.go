package service

import (
	"dungeon/internal/domain/models"
	"testing"
	"time"
)

func TestNewReporter(t *testing.T) {
	dungeon := createReporterTestDungeon()

	reporter := NewReporter(dungeon)

	if reporter == nil {
		t.Fatalf("reporter is nil")
	}

	if reporter.dungeon != dungeon {
		t.Fatalf("dungeon not assigned")
	}
}

func TestReporter_Generate(t *testing.T) {
	player := &models.Player{
		ID:             1,
		State:          models.StateCompleted,
		Health:         75,
		MonstersKilled: 40,
		BossKilled:     true,
		EnteredAt:      mustParseDateTime("2025-01-01 10:00:00"),
		LeftAt:         mustParseDateTime("2025-01-01 11:00:00"),
		FloorClearTimes: []time.Duration{
			5 * time.Minute,
			10 * time.Minute,
		},
		BossFloorEnteredAt: mustParseDateTime("2025-01-01 10:50:00"),
		BossKilledAt:       mustParseDateTime("2025-01-01 10:55:00"),
	}

	reporter := NewReporter(createReporterTestDungeon())

	report := reporter.Generate(player)

	if report.PlayerID != 1 {
		t.Errorf("expected player id 1")
	}

	if report.FinalHealth != 75 {
		t.Errorf("expected final health 75")
	}

	if report.State != models.FinalSuccess {
		t.Errorf("expected success state")
	}
}

func TestReporter_DetermineState(t *testing.T) {
	tests := []struct {
		name     string
		player   *models.Player
		expected models.FinalState
	}{
		{
			name: "disqualified",
			player: &models.Player{
				State: models.StateDisqualified,
			},
			expected: models.FinalDisqual,
		},
		{
			name: "dead",
			player: &models.Player{
				State: models.StateDead,
			},
			expected: models.FinalFail,
		},
		{
			name: "not completed and not in dungeon",
			player: &models.Player{
				State: models.StateRegistered,
			},
			expected: models.FinalDisqual,
		},
		{
			name: "successful completion",
			player: &models.Player{
				State:          models.StateCompleted,
				MonstersKilled: 40,
				BossKilled:     true,
			},
			expected: models.FinalSuccess,
		},
		{
			name: "boss not killed",
			player: &models.Player{
				State:          models.StateCompleted,
				MonstersKilled: 40,
				BossKilled:     false,
			},
			expected: models.FinalFail,
		},
		{
			name: "not enough monsters killed",
			player: &models.Player{
				State:          models.StateCompleted,
				MonstersKilled: 10,
				BossKilled:     true,
			},
			expected: models.FinalFail,
		},
	}

	reporter := NewReporter(createReporterTestDungeon())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reporter.determineState(tt.player)

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}

func TestReporter_CalculateTotalTime(t *testing.T) {
	reporter := NewReporter(createReporterTestDungeon())

	tests := []struct {
		name     string
		player   *models.Player
		expected time.Duration
	}{
		{
			name: "normal total time",
			player: &models.Player{
				EnteredAt: mustParseDateTime("2025-01-01 10:00:00"),
				LeftAt:    mustParseDateTime("2025-01-01 12:00:00"),
			},
			expected: 2 * time.Hour,
		},
		{
			name:     "zero entered time",
			player:   &models.Player{},
			expected: 0,
		},
		{
			name: "left before entered",
			player: &models.Player{
				EnteredAt: mustParseDateTime("2025-01-01 10:00:00"),
				LeftAt:    mustParseDateTime("2025-01-01 09:00:00"),
			},
			expected: 12 * time.Hour,
		},
		{
			name: "zero left time",
			player: &models.Player{
				EnteredAt: mustParseDateTime("2025-01-01 10:00:00"),
			},
			expected: 12 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reporter.calculateTotalTime(tt.player)

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}

func TestReporter_CalculateAvgFloorTime(t *testing.T) {
	reporter := NewReporter(createReporterTestDungeon())

	tests := []struct {
		name     string
		player   *models.Player
		expected time.Duration
	}{
		{
			name: "average calculated",
			player: &models.Player{
				FloorClearTimes: []time.Duration{
					5 * time.Minute,
					15 * time.Minute,
				},
			},
			expected: 10 * time.Minute,
		},
		{
			name:     "empty floor times",
			player:   &models.Player{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reporter.calculateAvgFloorTime(tt.player)

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}

func TestReporter_CalculateBossTime(t *testing.T) {
	reporter := NewReporter(createReporterTestDungeon())

	tests := []struct {
		name     string
		player   *models.Player
		expected time.Duration
	}{
		{
			name: "normal boss time",
			player: &models.Player{
				BossFloorEnteredAt: mustParseDateTime("2025-01-01 10:00:00"),
				BossKilledAt:       mustParseDateTime("2025-01-01 10:05:00"),
			},
			expected: 5 * time.Minute,
		},
		{
			name: "missing entered time",
			player: &models.Player{
				BossKilledAt: mustParseDateTime("2025-01-01 10:05:00"),
			},
			expected: 0,
		},
		{
			name: "missing killed time",
			player: &models.Player{
				BossFloorEnteredAt: mustParseDateTime("2025-01-01 10:00:00"),
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reporter.calculateBossTime(tt.player)

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}

func createReporterTestDungeon() *models.Dungeon {
	openTime := mustParseDateTime("2025-01-01 10:00:00")
	closeTime := mustParseDateTime("2025-01-01 22:00:00")

	return &models.Dungeon{
		Floors:           5,
		MonstersPerFloor: 10,
		OpenAt:           openTime,
		CloseAt:          closeTime,
	}
}

func mustParseDateTime(value string) time.Time {
	t, err := time.Parse(
		"2006-01-02 15:04:05",
		value,
	)
	if err != nil {
		panic(err)
	}
	return t
}
