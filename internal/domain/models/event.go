package models

import "time"

type Event struct {
	Time       time.Time
	PlayerID   int
	Type       EventType
	ExtraParam string
}

func NewEvent(t time.Time, playerID int, eventType EventType, extra string) *Event {
	return &Event{
		Time:       t,
		PlayerID:   playerID,
		Type:       eventType,
		ExtraParam: extra,
	}
}

func (e *Event) HasExtraParam() bool {
	return e.ExtraParam != ""
}
