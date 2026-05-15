package models

import (
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer(10)

	if player.ID != 10 {
		t.Errorf("expected ID 10, got %d", player.ID)
	}

	if player.State != StateUnregistered {
		t.Errorf("unexpected initial state")
	}

	if player.Health != StartHealth {
		t.Errorf("expected health %d, got %d",
			StartHealth,
			player.Health)
	}

	if player.CurrentFloor != StartFloor {
		t.Errorf("expected floor %d, got %d",
			StartFloor,
			player.CurrentFloor)
	}

	if len(player.FloorClearTimes) != 0 {
		t.Errorf("expected empty floor clear times")
	}
}

func TestPlayer_Register(t *testing.T) {
	player := NewPlayer(1)

	player.Register()

	if player.State != StateRegistered {
		t.Errorf("player should be registered")
	}
}

func TestPlayer_EnterDungeon(t *testing.T) {
	player := NewPlayer(1)

	now := time.Now()

	player.EnterDungeon(now)

	if player.State != StateInDungeon {
		t.Errorf("expected in dungeon state")
	}

	if player.EnteredAt != now {
		t.Errorf("EnteredAt not set")
	}

	if player.CurrentFloor != 1 {
		t.Errorf("expected floor 1")
	}
}

func TestPlayer_LeaveDungeon(t *testing.T) {
	player := NewPlayer(1)
	player.State = StateInDungeon

	now := time.Now()

	player.LeaveDungeon(now)

	if player.LeftAt != now {
		t.Errorf("LeftAt not set")
	}

	if player.State != StateCompleted {
		t.Errorf("expected completed state")
	}
}

func TestPlayer_MoveToNextFloor(t *testing.T) {
	player := NewPlayer(1)
	player.CurrentFloor = 1

	ok := player.MoveToNextFloor(time.Now(), 5)

	if !ok {
		t.Errorf("expected move success")
	}

	if player.CurrentFloor != 2 {
		t.Errorf("expected floor 2, got %d",
			player.CurrentFloor)
	}
}

func TestPlayer_MoveToNextFloor_MaxFloor(t *testing.T) {
	player := NewPlayer(1)
	player.CurrentFloor = 5

	ok := player.MoveToNextFloor(time.Now(), 5)

	if ok {
		t.Errorf("expected move failure")
	}

	if player.CurrentFloor != 5 {
		t.Errorf("floor should not change")
	}
}

func TestPlayer_MoveToPreviousFloor(t *testing.T) {
	player := NewPlayer(1)
	player.CurrentFloor = 3

	ok := player.MoveToPreviousFloor(time.Now())

	if !ok {
		t.Errorf("expected successful move")
	}

	if player.CurrentFloor != 2 {
		t.Errorf("expected floor 2")
	}
}

func TestPlayer_MoveToPreviousFloor_FirstFloor(t *testing.T) {
	player := NewPlayer(1)
	player.CurrentFloor = 1

	ok := player.MoveToPreviousFloor(time.Now())

	if ok {
		t.Errorf("expected failure")
	}

	if player.CurrentFloor != 1 {
		t.Errorf("floor should stay 1")
	}
}

func TestPlayer_KillMonster(t *testing.T) {
	player := NewPlayer(1)

	player.KillMonster()
	player.KillMonster()

	if player.MonstersKilled != 2 {
		t.Errorf("expected 2 monsters killed")
	}
}

func TestPlayer_KillBoss(t *testing.T) {
	player := NewPlayer(1)

	now := time.Now()

	player.KillBoss(now)

	if !player.BossKilled {
		t.Errorf("boss should be killed")
	}

	if player.BossKilledAt != now {
		t.Errorf("BossKilledAt not set")
	}
}

func TestPlayer_CompleteFloor(t *testing.T) {
	player := NewPlayer(1)

	start := time.Now()
	end := start.Add(5 * time.Minute)

	player.FloorEnteredAt = start

	player.CompleteFloor(end)

	if len(player.FloorClearTimes) != 1 {
		t.Errorf("expected 1 floor clear time")
	}

	if player.FloorClearTimes[0] != 5*time.Minute {
		t.Errorf("unexpected clear time")
	}
}

func TestPlayer_TakeDamage(t *testing.T) {
	tests := []struct {
		name         string
		startHealth  int
		damage       int
		expectedHP   int
		expectedDead bool
	}{
		{
			name:         "normal damage",
			startHealth:  100,
			damage:       30,
			expectedHP:   70,
			expectedDead: false,
		},
		{
			name:         "fatal damage",
			startHealth:  50,
			damage:       100,
			expectedHP:   0,
			expectedDead: true,
		},
		{
			name:         "exact death",
			startHealth:  40,
			damage:       40,
			expectedHP:   0,
			expectedDead: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1)
			player.Health = tt.startHealth

			player.TakeDamage(tt.damage)

			if player.Health != tt.expectedHP {
				t.Errorf("expected hp %d, got %d",
					tt.expectedHP,
					player.Health)
			}

			if tt.expectedDead && player.State != StateDead {
				t.Errorf("expected dead state")
			}
		})
	}
}

func TestPlayer_Heal(t *testing.T) {
	tests := []struct {
		name       string
		startHP    int
		heal       int
		expectedHP int
	}{
		{
			name:       "normal heal",
			startHP:    50,
			heal:       20,
			expectedHP: 70,
		},
		{
			name:       "overheal",
			startHP:    90,
			heal:       50,
			expectedHP: MaxHealth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1)
			player.Health = tt.startHP

			player.Heal(tt.heal)

			if player.Health != tt.expectedHP {
				t.Errorf("expected hp %d, got %d",
					tt.expectedHP,
					player.Health)
			}
		})
	}
}

func TestPlayer_Disqualify(t *testing.T) {
	player := NewPlayer(1)

	player.Disqualify()

	if player.State != StateDisqualified {
		t.Errorf("expected disqualified state")
	}
}

func TestPlayer_IsAlive(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		want   bool
	}{
		{
			name: "alive",
			player: Player{
				Health: 10,
				State:  StateInDungeon,
			},
			want: true,
		},
		{
			name: "dead by health",
			player: Player{
				Health: 0,
				State:  StateInDungeon,
			},
			want: false,
		},
		{
			name: "dead state",
			player: Player{
				Health: 10,
				State:  StateDead,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.player.IsAlive()

			if got != tt.want {
				t.Errorf("expected %v, got %v",
					tt.want,
					got)
			}
		})
	}
}

func TestPlayer_CanParticipate(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		want   bool
	}{
		{
			name: "registered",
			player: Player{
				State:  StateRegistered,
				Health: 100,
			},
			want: true,
		},
		{
			name: "in dungeon",
			player: Player{
				State:  StateInDungeon,
				Health: 100,
			},
			want: true,
		},
		{
			name: "dead",
			player: Player{
				State:  StateInDungeon,
				Health: 0,
			},
			want: false,
		},
		{
			name: "disqualified",
			player: Player{
				State:  StateDisqualified,
				Health: 100,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.player.CanParticipate()

			if got != tt.want {
				t.Errorf("expected %v, got %v",
					tt.want,
					got)
			}
		})
	}
}
