package models

import "testing"

func TestEventType_String(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		expected  string
	}{
		{
			name:      "player registered",
			eventType: EventPlayerRegistered,
			expected:  "PlayerRegistered",
		},
		{
			name:      "player entered",
			eventType: EventPlayerEntered,
			expected:  "PlayerEntered",
		},
		{
			name:      "monster killed",
			eventType: EventMonsterKilled,
			expected:  "MonsterKilled",
		},
		{
			name:      "boss killed",
			eventType: EventBossKilled,
			expected:  "BossKilled",
		},
		{
			name:      "unknown event",
			eventType: EventType(999),
			expected:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.eventType.String()

			if got != tt.expected {
				t.Errorf("expected %s, got %s",
					tt.expected,
					got)
			}
		})
	}
}

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		expected  bool
	}{
		{
			name:      "valid first event",
			eventType: EventPlayerRegistered,
			expected:  true,
		},
		{
			name:      "valid last event",
			eventType: EventDamageReceived,
			expected:  true,
		},
		{
			name:      "below valid range",
			eventType: EventType(0),
			expected:  false,
		},
		{
			name:      "above valid range",
			eventType: EventType(999),
			expected:  false,
		},
		{
			name:      "negative event type",
			eventType: EventType(-1),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.eventType.IsValid()

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if MaxHealth != 100 {
		t.Errorf("expected MaxHealth = 100")
	}

	if StartHealth != 100 {
		t.Errorf("expected StartHealth = 100")
	}

	if StartFloor != 0 {
		t.Errorf("expected StartFloor = 0")
	}
}

func TestPlayerStates(t *testing.T) {
	tests := []struct {
		state    PlayerState
		expected string
	}{
		{StateUnregistered, "UNREGISTERED"},
		{StateRegistered, "REGISTERED"},
		{StateInDungeon, "IN_DUNGEON"},
		{StateCompleted, "COMPLETED"},
		{StateDead, "DEAD"},
		{StateDisqualified, "DISQUALIFIED"},
	}

	for _, tt := range tests {
		if string(tt.state) != tt.expected {
			t.Errorf("expected %s, got %s",
				tt.expected,
				tt.state)
		}
	}
}

func TestFinalStates(t *testing.T) {
	tests := []struct {
		state    FinalState
		expected string
	}{
		{FinalSuccess, "SUCCESS"},
		{FinalFail, "FAIL"},
		{FinalDisqual, "DISQUAL"},
	}

	for _, tt := range tests {
		if string(tt.state) != tt.expected {
			t.Errorf("expected %s, got %s",
				tt.expected,
				tt.state)
		}
	}
}
