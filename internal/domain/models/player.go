package models

import (
	"time"
)

type Player struct {
	ID             int
	State          PlayerState
	Health         int
	CurrentFloor   int
	MonstersKilled int
	BossKilled     bool

	EnteredAt          time.Time
	LeftAt             time.Time
	FloorEnteredAt     time.Time
	BossFloorEnteredAt time.Time
	BossKilledAt       time.Time

	FloorClearTimes []time.Duration
}

func NewPlayer(id int) *Player {
	return &Player{
		ID:              id,
		State:           StateUnregistered,
		Health:          StartHealth,
		CurrentFloor:    StartFloor,
		FloorClearTimes: make([]time.Duration, 0),
	}
}

// Register transitions the player to the registered state if they are currently unregistered
func (p *Player) Register() {
	if p.State == StateUnregistered {
		p.State = StateRegistered
	}
}

// EnterDungeon transitions the player to the in-dungeon state and records the entry time
func (p *Player) EnterDungeon(t time.Time) {
	p.State = StateInDungeon
	p.EnteredAt = t
	p.FloorEnteredAt = t
	p.CurrentFloor = 1
}

// LeaveDungeon transitions the player to the completed state if they are currently in the dungeon and records the exit time
func (p *Player) LeaveDungeon(t time.Time) {
	p.LeftAt = t
	if p.State == StateInDungeon {
		p.State = StateCompleted
	}
}

// MoveToNextFloor attempts to move the player to the next floor, updating the current floor and entry time if successful
func (p *Player) MoveToNextFloor(t time.Time, maxFloors int) bool {
	if p.CurrentFloor >= maxFloors {
		return false
	}
	p.CurrentFloor++
	p.FloorEnteredAt = t
	return true
}

// MoveToPreviousFloor attempts to move the player to the previous floor, updating the current floor and entry time if successful
func (p *Player) MoveToPreviousFloor(t time.Time) bool {
	if p.CurrentFloor <= 1 {
		return false
	}
	p.CurrentFloor--
	p.FloorEnteredAt = t
	return true
}

// EnterBossFloor transitions the player to the boss floor state and records the entry time
func (p *Player) EnterBossFloor(t time.Time) {
	p.BossFloorEnteredAt = t
}

// KillMonster increments the count of monsters killed by the player
func (p *Player) KillMonster() {
	p.MonstersKilled++
}

// KillBoss transitions the player to the boss killed state and records the time of the kill
func (p *Player) KillBoss(t time.Time) {
	p.BossKilled = true
	p.BossKilledAt = t
}

// CompleteFloor records the time taken to clear the current floor and adds it to the player's floor clear times
func (p *Player) CompleteFloor(t time.Time) {
	clearTime := t.Sub(p.FloorEnteredAt)
	p.FloorClearTimes = append(p.FloorClearTimes, clearTime)
}

// TakeDamage reduces the player's health by the specified damage amount and updates the player's state if health drops to zero
func (p *Player) TakeDamage(damage int) {
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}
	if p.Health == 0 {
		p.State = StateDead
	}
}

// Heal increases the player's health by the specified amount, ensuring it does not exceed the maximum health
func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > MaxHealth {
		p.Health = MaxHealth
	}
}

// Disqualify transitions the player to the disqualified state, indicating they can no longer participate in the dungeon
func (p *Player) Disqualify() {
	p.State = StateDisqualified
}

// IsAlive checks if the player is alive based on their health and state
func (p *Player) IsAlive() bool {
	return p.Health > 0 && p.State != StateDead
}

// IsRegistered checks if the player is registered for the dungeon
func (p *Player) IsRegistered() bool {
	return p.State != StateUnregistered
}

// IsInDungeon checks if the player is currently in the dungeon
func (p *Player) IsInDungeon() bool {
	return p.State == StateInDungeon
}

// CanParticipate checks if the player is in a state that allows them to participate in the dungeon (registered or in-dungeon) and is alive
func (p *Player) CanParticipate() bool {
	return (p.State == StateRegistered || p.State == StateInDungeon) && p.IsAlive()
}
