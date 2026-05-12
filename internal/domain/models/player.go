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

func (p *Player) Register() {
	if p.State == StateUnregistered {
		p.State = StateRegistered
	}
}

func (p *Player) EnterDungeon(t time.Time) {
	p.State = StateInDungeon
	p.EnteredAt = t
	p.FloorEnteredAt = t
	p.CurrentFloor = 1
}

func (p *Player) LeaveDungeon(t time.Time) {
	p.LeftAt = t
	if p.State == StateInDungeon {
		p.State = StateCompleted
	}
}

func (p *Player) MoveToNextFloor(t time.Time, maxFloors int) bool {
	if p.CurrentFloor >= maxFloors {
		return false
	}
	p.CurrentFloor++
	p.FloorEnteredAt = t
	return true
}

func (p *Player) MoveToPreviousFloor(t time.Time) bool {
	if p.CurrentFloor <= 1 {
		return false
	}
	p.CurrentFloor--
	p.FloorEnteredAt = t
	return true
}

func (p *Player) EnterBossFloor(t time.Time) {
	p.BossFloorEnteredAt = t
}

func (p *Player) KillMonster() {
	p.MonstersKilled++
}

func (p *Player) KillBoss(t time.Time) {
	p.BossKilled = true
	p.BossKilledAt = t
}

func (p *Player) CompleteFloor(t time.Time) {
	clearTime := t.Sub(p.FloorEnteredAt)
	p.FloorClearTimes = append(p.FloorClearTimes, clearTime)
}

func (p *Player) TakeDamage(damage int) {
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}
	if p.Health == 0 {
		p.State = StateDead
	}
}

func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > MaxHealth {
		p.Health = MaxHealth
	}
}

func (p *Player) Disqualify() {
	p.State = StateDisqualified
}

func (p *Player) IsAlive() bool {
	return p.Health > 0 && p.State != StateDead
}

func (p *Player) IsRegistered() bool {
	return p.State != StateUnregistered
}

func (p *Player) IsInDungeon() bool {
	return p.State == StateInDungeon
}

func (p *Player) CanParticipate() bool {
	return (p.State == StateRegistered || p.State == StateInDungeon) && p.IsAlive()
}
