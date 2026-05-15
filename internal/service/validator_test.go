package service

import (
	"dungeon/internal/domain/models"
	"testing"
	"time"
)

func TestNewValidator(t *testing.T) {
	dungeon := createTestDungeon()

	validator := NewValidator(dungeon)

	if validator == nil {
		t.Fatalf("validator is nil")
	}

	if validator.dungeon != dungeon {
		t.Fatalf("dungeon not assigned")
	}
}

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name        string
		event       *models.Event
		player      *models.Player
		expectError bool
	}{
		{
			name: "unregistered player action",
			event: createEvent(
				"12:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:    1,
				State: models.StateUnregistered,
			},
			expectError: true,
		},
		{
			name: "registered player action",
			event: createEvent(
				"12:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 1,
			},
			expectError: false,
		},
		{
			name: "dungeon closed",
			event: createEvent(
				"23:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:     1,
				State:  models.StateInDungeon,
				Health: 100,
			},
			expectError: false,
		},
		{
			name: "dead player ignored",
			event: createEvent(
				"12:00:00",
				models.EventNextFloor,
			),
			player: &models.Player{
				ID:     1,
				State:  models.StateDead,
				Health: 0,
			},
			expectError: false,
		},
		{
			name: "impossible previous floor move",
			event: createEvent(
				"12:00:00",
				models.EventPreviousFloor,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 1,
			},
			expectError: true,
		},
		{
			name: "valid previous floor move",
			event: createEvent(
				"12:00:00",
				models.EventPreviousFloor,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 3,
			},
			expectError: false,
		},
		{
			name: "impossible next floor move",
			event: createEvent(
				"12:00:00",
				models.EventNextFloor,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 5,
			},
			expectError: true,
		},
		{
			name: "valid next floor move",
			event: createEvent(
				"12:00:00",
				models.EventNextFloor,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 2,
			},
			expectError: false,
		},
		{
			name: "enter dungeon twice",
			event: createEvent(
				"12:00:00",
				models.EventPlayerEntered,
			),
			player: &models.Player{
				ID:     1,
				State:  models.StateInDungeon,
				Health: 100,
			},
			expectError: true,
		},
		{
			name: "monster kill outside dungeon",
			event: createEvent(
				"12:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:     1,
				State:  models.StateRegistered,
				Health: 100,
			},
			expectError: true,
		},
		{
			name: "monster kill on boss floor",
			event: createEvent(
				"12:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 5,
			},
			expectError: true,
		},
		{
			name: "valid monster kill",
			event: createEvent(
				"12:00:00",
				models.EventMonsterKilled,
			),
			player: &models.Player{
				ID:           1,
				State:        models.StateInDungeon,
				Health:       100,
				CurrentFloor: 2,
			},
			expectError: false,
		},
	}

	validator := NewValidator(createTestDungeon())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.event, tt.player)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func createTestDungeon() *models.Dungeon {
	openTime, _ := time.Parse("15:04:05", "10:00:00")
	closeTime, _ := time.Parse("15:04:05", "22:00:00")

	return &models.Dungeon{
		Floors:           5,
		MonstersPerFloor: 10,
		OpenAt:           openTime,
		CloseAt:          closeTime,
	}
}

func createEvent(
	timeStr string,
	eventType models.EventType,
) *models.Event {

	t, _ := time.Parse("15:04:05", timeStr)

	return &models.Event{
		Time:     t,
		PlayerID: 1,
		Type:     eventType,
	}
}
