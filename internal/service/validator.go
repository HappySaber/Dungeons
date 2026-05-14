package service

import (
	"dungeon/internal/domain/models"
	"fmt"
)

type Validator struct {
	dungeon *models.Dungeon
}

func NewValidator(dungeon *models.Dungeon) *Validator {
	return &Validator{
		dungeon: dungeon,
	}
}

func (v *Validator) Validate(event *models.Event, player *models.Player) error {
	t := event.Time

	if event.Type != models.EventPlayerRegistered && !player.IsRegistered() {
		player.Disqualify()
		return fmt.Errorf("[%s] Player [%d] is disqualified", formatTime(t), player.ID)
	}

	if event.Type != models.EventPlayerRegistered && !v.dungeon.IsOpen(t) {

		return nil
	}

	if !player.IsAlive() && event.Type != models.EventDamageReceived {
		return nil
	}

	switch event.Type {
	case models.EventPreviousFloor:
		if player.CurrentFloor <= 1 {
			return fmt.Errorf("[%s] Player [%d] makes imposible move [%d]",
				formatTime(t), player.ID, event.Type)
		}

	case models.EventNextFloor:
		if player.CurrentFloor >= v.dungeon.Floors {
			return fmt.Errorf("[%s] Player [%d] makes imposible move [%d]",
				formatTime(t), player.ID, event.Type)
		}

	case models.EventPlayerEntered:
		if player.IsInDungeon() {
			return fmt.Errorf("[%s] Player [%d] makes imposible move [%d]",
				formatTime(t), player.ID, event.Type)
		}

	case models.EventMonsterKilled:
		if !player.IsInDungeon() {
			return fmt.Errorf("[%s] Player [%d] makes imposible move [%d]",
				formatTime(t), player.ID, event.Type)
		}
		if v.dungeon.IsBossFloor(player.CurrentFloor) {
			return fmt.Errorf("[%s] Player [%d] makes imposible move [%d]",
				formatTime(t), player.ID, event.Type)
		}
	}

	return nil
}
