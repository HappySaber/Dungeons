package parser

import (
	"dungeon/internal/domain/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EventParser struct {
	regex *regexp.Regexp
}

func New() *EventParser {
	pattern := `^\[(\d{2}:\d{2}:\d{2})\]\s+(\d+)\s+(\d+)(?:\s+(.+))?$`

	return &EventParser{
		regex: regexp.MustCompile(pattern),
	}
}

func (p *EventParser) Parse(line string) (*models.Event, error) {
	line = strings.TrimSpace(line)

	matches := p.regex.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid event format: %s", line)
	}

	eventTime, err := time.Parse("15:04:05", matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid time format %q: %w", matches[1], err)
	}

	playerID, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid player ID %q: %w", matches[2], err)
	}

	eventTypeNum, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid event type %q: %w", matches[3], err)
	}

	eventType := models.EventType(eventTypeNum)
	if !eventType.IsValid() {
		return nil, fmt.Errorf("unknown event type: %d", eventTypeNum)
	}

	extraParam := ""
	if len(matches) > 4 && matches[4] != "" {
		extraParam = strings.TrimSpace(matches[4])
	}

	return models.NewEvent(eventTime, playerID, eventType, extraParam), nil
}
