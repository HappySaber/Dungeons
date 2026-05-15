package timeutil

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00:00",
		},
		{
			name:     "seconds only",
			duration: 45 * time.Second,
			expected: "00:00:45",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "00:05:30",
		},
		{
			name:     "hours minutes seconds",
			duration: 2*time.Hour + 15*time.Minute + 10*time.Second,
			expected: "02:15:10",
		},
		{
			name:     "more than 24 hours",
			duration: 27*time.Hour + 5*time.Minute,
			expected: "27:05:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.duration)

			if got != tt.expected {
				t.Errorf("expected %s, got %s",
					tt.expected,
					got)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid time",
			input:       "12:30:45",
			expectError: false,
		},
		{
			name:        "midnight",
			input:       "00:00:00",
			expectError: false,
		},
		{
			name:        "invalid hour",
			input:       "25:00:00",
			expectError: true,
		},
		{
			name:        "invalid minute",
			input:       "12:99:00",
			expectError: true,
		},
		{
			name:        "invalid format",
			input:       "12-30-45",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseTime(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			formatted := parsed.Format("15:04:05")

			if formatted != tt.input {
				t.Errorf("expected %s, got %s",
					tt.input,
					formatted)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "zero duration",
			input:       "00:00:00",
			expected:    0,
			expectError: false,
		},
		{
			name:        "seconds only",
			input:       "00:00:30",
			expected:    30 * time.Second,
			expectError: false,
		},
		{
			name:        "minutes and seconds",
			input:       "00:05:15",
			expected:    5*time.Minute + 15*time.Second,
			expectError: false,
		},
		{
			name:        "hours minutes seconds",
			input:       "02:10:45",
			expected:    2*time.Hour + 10*time.Minute + 45*time.Second,
			expectError: false,
		},
		{
			name:        "invalid format",
			input:       "2h10m",
			expectError: true,
		},
		{
			name:        "invalid time",
			input:       "99:99:99",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Errorf("expected %v, got %v",
					tt.expected,
					got)
			}
		})
	}
}
