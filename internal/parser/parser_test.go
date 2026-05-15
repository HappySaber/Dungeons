package parser

import (
	"dungeon/internal/domain/models"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	parser := New()

	if parser == nil {
		t.Fatalf("expected parser, got nil")
	}

	if parser.regex == nil {
		t.Fatalf("regex should not be nil")
	}
}

func TestEventParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		expected    *models.Event
		expectError bool
	}{
		{
			name: "valid event without extra param",
			line: "[12:30:45] 42 1",
			expected: &models.Event{
				Time:       mustParseTime("12:30:45"),
				PlayerID:   42,
				Type:       models.EventPlayerRegistered,
				ExtraParam: "",
			},
			expectError: false,
		},
		{
			name: "valid event with extra param",
			line: "[08:15:00] 7 3 dragon",
			expected: &models.Event{
				Time:       mustParseTime("08:15:00"),
				PlayerID:   7,
				Type:       models.EventMonsterKilled,
				ExtraParam: "dragon",
			},
			expectError: false,
		},
		{
			name: "valid event with spaces",
			line: "   [10:00:00] 1 2   boss room   ",
			expected: &models.Event{
				Time:       mustParseTime("10:00:00"),
				PlayerID:   1,
				Type:       models.EventPlayerEntered,
				ExtraParam: "boss room",
			},
			expectError: false,
		},
		{
			name:        "invalid format",
			line:        "wrong format",
			expectError: true,
		},
		{
			name:        "invalid time",
			line:        "[99:99:99] 1 1",
			expectError: true,
		},
		{
			name:        "invalid player id",
			line:        "[12:00:00] abc 1",
			expectError: true,
		},
		{
			name:        "invalid event type",
			line:        "[12:00:00] 1 abc",
			expectError: true,
		},
		{
			name:        "unknown event type",
			line:        "[12:00:00] 1 999",
			expectError: true,
		},
		{
			name:        "missing fields",
			line:        "[12:00:00] 1",
			expectError: true,
		},
		{
			name:        "empty line",
			line:        "",
			expectError: true,
		},
	}

	parser := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := parser.Parse(tt.line)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if event.PlayerID != tt.expected.PlayerID {
				t.Errorf("expected player id %d, got %d",
					tt.expected.PlayerID,
					event.PlayerID)
			}

			if event.Type != tt.expected.Type {
				t.Errorf("expected event type %v, got %v",
					tt.expected.Type,
					event.Type)
			}

			if event.ExtraParam != tt.expected.ExtraParam {
				t.Errorf("expected extra param %q, got %q",
					tt.expected.ExtraParam,
					event.ExtraParam)
			}

			if !event.Time.Equal(tt.expected.Time) {
				t.Errorf("expected time %v, got %v",
					tt.expected.Time,
					event.Time)
			}
		})
	}
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse("15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}
