package models

import (
	"dungeon/internal/config"
	"testing"
	"time"
)

func TestParseOpenTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid time",
			input:   "12:30:00",
			wantErr: false,
		},
		{
			name:    "invalid time",
			input:   "99:99:99",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseOpenTime(tt.input)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCalculateCloseTime(t *testing.T) {
	openTime := time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)

	closeTime := CalculateCloseTime(5, openTime)

	expected := openTime.Add(5 * time.Hour)

	if !closeTime.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, closeTime)
	}
}

func TestNewDungeonFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config.Config{
				Floors:   10,
				Monsters: 5,
				OpenAt:   "08:00:00",
				Duration: 6,
			},
			wantErr: false,
		},
		{
			name: "invalid open time",
			cfg: &config.Config{
				Floors:   10,
				Monsters: 5,
				OpenAt:   "invalid",
				Duration: 6,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dungeon, err := NewDungeonFromConfig(tt.cfg)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr {
				if dungeon.Floors != tt.cfg.Floors {
					t.Errorf("expected Floors %d, got %d",
						tt.cfg.Floors,
						dungeon.Floors)
				}

				if dungeon.MonstersPerFloor != tt.cfg.Monsters {
					t.Errorf("expected MonstersPerFloor %d, got %d",
						tt.cfg.Monsters,
						dungeon.MonstersPerFloor)
				}
			}
		})
	}
}

func TestDungeon_IsOpen(t *testing.T) {
	open := time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)
	close := open.Add(2 * time.Hour)

	dungeon := &Dungeon{
		OpenAt:  open,
		CloseAt: close,
	}

	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		{
			name: "during working hours",
			time: open.Add(1 * time.Hour),
			want: true,
		},
		{
			name: "before open",
			time: open.Add(-1 * time.Hour),
			want: false,
		},
		{
			name: "after close",
			time: close.Add(1 * time.Minute),
			want: false,
		},
		{
			name: "exact open time",
			time: open,
			want: true,
		},
		{
			name: "exact close time",
			time: close,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dungeon.IsOpen(tt.time)

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestDungeon_IsBossFloor(t *testing.T) {
	dungeon := &Dungeon{
		Floors: 5,
	}

	if !dungeon.IsBossFloor(5) {
		t.Errorf("floor 5 should be boss floor")
	}

	if dungeon.IsBossFloor(4) {
		t.Errorf("floor 4 should not be boss floor")
	}
}

func TestDungeon_IsValidFloor(t *testing.T) {
	dungeon := &Dungeon{
		Floors: 5,
	}

	tests := []struct {
		floor int
		want  bool
	}{
		{1, true},
		{5, true},
		{0, false},
		{6, false},
		{-1, false},
	}

	for _, tt := range tests {
		got := dungeon.IsValidFloor(tt.floor)

		if got != tt.want {
			t.Errorf("floor %d: expected %v, got %v",
				tt.floor,
				tt.want,
				got)
		}
	}
}

func TestDungeon_TotalMonsters(t *testing.T) {
	dungeon := &Dungeon{
		Floors:           5,
		MonstersPerFloor: 10,
	}

	got := dungeon.TotalMonsters()
	want := 40

	if got != want {
		t.Errorf("expected %d, got %d", want, got)
	}
}

func TestDungeon_MonstersOnFloor(t *testing.T) {
	dungeon := &Dungeon{
		Floors:           5,
		MonstersPerFloor: 10,
	}

	tests := []struct {
		floor int
		want  int
	}{
		{1, 10},
		{2, 10},
		{5, 0},
	}

	for _, tt := range tests {
		got := dungeon.MonstersOnFloor(tt.floor)

		if got != tt.want {
			t.Errorf("floor %d: expected %d, got %d",
				tt.floor,
				tt.want,
				got)
		}
	}
}
