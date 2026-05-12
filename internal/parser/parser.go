package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dungeon/internal/domain/models"
)

const clockLayout = "15:04:05"

func ParseEvents(path string) ([]models.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []models.Event

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		event, err := ParseEventLine(line)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, scanner.Err()
}

func ParseEventLine(line string) (models.Event, error) {
	parts := strings.Split(line, " ")

	if len(parts) < 3 {
		return models.Event{}, fmt.Errorf("invalid event")
	}

	timeStr := strings.Trim(parts[0], "[]")

	t, err := ParseClock(timeStr)
	if err != nil {
		return models.Event{}, err
	}

	playerID, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Event{}, err
	}

	eventID, err := strconv.Atoi(parts[2])
	if err != nil {
		return models.Event{}, err
	}

	extra := ""

	if len(parts) > 3 {
		extra = strings.Join(parts[3:], " ")
	}

	return models.Event{
		Time:       t,
		PlayerID:   playerID,
		ID:         eventID,
		ExtraParam: extra,
	}, nil
}

func ParseClock(value string) (time.Time, error) {
	return time.Parse(clockLayout, value)
}

func FormatClock(t time.Time) string {
	return t.Format(clockLayout)
}
