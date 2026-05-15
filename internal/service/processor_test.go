package service

import (
	"dungeon/internal/domain/models"
	"testing"
	"time"
)

func TestNewProcessor(t *testing.T) {
	dungeon := createProcessorTestDungeon()

	processor := NewProcessor(dungeon)

	if processor == nil {
		t.Fatalf("processor is nil")
	}

	if processor.dungeon != dungeon {
		t.Fatalf("dungeon not assigned")
	}

	if processor.players == nil {
		t.Fatalf("players map not initialized")
	}

	if processor.validator == nil {
		t.Fatalf("validator not initialized")
	}

	if processor.reporter == nil {
		t.Fatalf("reporter not initialized")
	}
}

func TestProcessor_GetOrCreatePlayer(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	player1 := processor.getOrCreatePlayer(1)
	player2 := processor.getOrCreatePlayer(1)

	if player1 != player2 {
		t.Fatalf("expected same player instance")
	}

	if len(processor.players) != 1 {
		t.Fatalf("expected 1 player")
	}
}

func TestProcessor_ProcessRegisterEvent(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	event := createProcessorEvent(
		"12:00:00",
		1,
		models.EventPlayerRegistered,
		"",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player == nil {
		t.Fatalf("player not created")
	}

	if player.State != models.StateRegistered {
		t.Fatalf("player should be registered")
	}

	if len(processor.output) != 1 {
		t.Fatalf("expected 1 output message")
	}
}

func TestProcessor_ProcessEnterDungeon(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	register := createProcessorEvent(
		"12:00:00",
		1,
		models.EventPlayerRegistered,
		"",
	)

	enter := createProcessorEvent(
		"12:01:00",
		1,
		models.EventPlayerEntered,
		"",
	)

	processor.ProcessEvent(register)
	processor.ProcessEvent(enter)

	player := processor.players[1]

	if !player.IsInDungeon() {
		t.Fatalf("player should be in dungeon")
	}
}

func TestProcessor_ProcessMonsterKill(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	event := createProcessorEvent(
		"12:05:00",
		1,
		models.EventMonsterKilled,
		"",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.MonstersKilled != 1 {
		t.Fatalf("expected 1 killed monster")
	}
}

func TestProcessor_ProcessNextFloor(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	event := createProcessorEvent(
		"12:10:00",
		1,
		models.EventNextFloor,
		"",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.CurrentFloor != 2 {
		t.Fatalf("expected floor 2")
	}
}

func TestProcessor_ProcessPreviousFloor(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	player := processor.players[1]
	player.CurrentFloor = 2

	event := createProcessorEvent(
		"12:10:00",
		1,
		models.EventPreviousFloor,
		"",
	)

	processor.ProcessEvent(event)

	if player.CurrentFloor != 1 {
		t.Fatalf("expected floor 1")
	}
}

func TestProcessor_ProcessBossKill(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	player := processor.players[1]
	player.CurrentFloor = 5

	enterBoss := createProcessorEvent(
		"12:20:00",
		1,
		models.EventBossFloor,
		"",
	)

	killBoss := createProcessorEvent(
		"12:25:00",
		1,
		models.EventBossKilled,
		"",
	)

	processor.ProcessEvent(enterBoss)
	processor.ProcessEvent(killBoss)

	if !player.BossKilled {
		t.Fatalf("boss should be killed")
	}
}

func TestProcessor_ProcessDamage(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	event := createProcessorEvent(
		"12:30:00",
		1,
		models.EventDamageReceived,
		"30",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.Health != 70 {
		t.Fatalf("expected health 70, got %d",
			player.Health)
	}
}

func TestProcessor_ProcessFatalDamage(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	event := createProcessorEvent(
		"12:30:00",
		1,
		models.EventDamageReceived,
		"200",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.State != models.StateDead {
		t.Fatalf("player should be dead")
	}

	if len(processor.output) < 2 {
		t.Fatalf("expected death output message")
	}
}

func TestProcessor_ProcessHeal(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	player := processor.players[1]
	player.Health = 50

	event := createProcessorEvent(
		"12:35:00",
		1,
		models.EventHealthRestored,
		"20",
	)

	processor.ProcessEvent(event)

	if player.Health != 70 {
		t.Fatalf("expected health 70")
	}
}

func TestProcessor_ProcessDisqualify(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	runBasicFlow(processor)

	event := createProcessorEvent(
		"12:40:00",
		1,
		models.EventCannotContinue,
		"panic",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.State != models.StateDisqualified {
		t.Fatalf("expected disqualified state")
	}
}

func TestProcessor_InvalidAction(t *testing.T) {
	processor := NewProcessor(createProcessorTestDungeon())

	event := createProcessorEvent(
		"12:00:00",
		1,
		models.EventMonsterKilled,
		"",
	)

	processor.ProcessEvent(event)

	player := processor.players[1]

	if player.State != models.StateDisqualified {
		t.Fatalf("player should be disqualified")
	}

	if len(processor.output) == 0 {
		t.Fatalf("expected validation output")
	}
}

func TestParseIntParam(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"10", 10},
		{"0", 0},
		{"999", 999},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		got := parseIntParam(tt.input)

		if got != tt.expected {
			t.Fatalf("expected %d, got %d",
				tt.expected,
				got)
		}
	}
}

func TestFormatTime(t *testing.T) {
	tm := time.Date(
		2025,
		1,
		1,
		12,
		34,
		56,
		0,
		time.UTC,
	)

	got := formatTime(tm)

	if got != "12:34:56" {
		t.Fatalf("unexpected formatted time: %s", got)
	}
}

func runBasicFlow(processor *Processor) {
	register := createProcessorEvent(
		"12:00:00",
		1,
		models.EventPlayerRegistered,
		"",
	)

	enter := createProcessorEvent(
		"12:01:00",
		1,
		models.EventPlayerEntered,
		"",
	)

	processor.ProcessEvent(register)
	processor.ProcessEvent(enter)
}

func createProcessorTestDungeon() *models.Dungeon {
	openTime, _ := time.Parse("15:04:05", "10:00:00")
	closeTime, _ := time.Parse("15:04:05", "22:00:00")

	return &models.Dungeon{
		Floors:           5,
		MonstersPerFloor: 10,
		OpenAt:           openTime,
		CloseAt:          closeTime,
	}
}

func createProcessorEvent(
	timeStr string,
	playerID int,
	eventType models.EventType,
	extra string,
) *models.Event {

	tm, _ := time.Parse("15:04:05", timeStr)

	return &models.Event{
		Time:       tm,
		PlayerID:   playerID,
		Type:       eventType,
		ExtraParam: extra,
	}
}
