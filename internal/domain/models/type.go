package models

type PlayerState string

const (
	StateUnregistered PlayerState = "UNREGISTERED"
	StateRegistered   PlayerState = "REGISTERED"
	StateInDungeon    PlayerState = "IN_DUNGEON"
	StateCompleted    PlayerState = "COMPLETED"
	StateDead         PlayerState = "DEAD"
	StateDisqualified PlayerState = "DISQUALIFIED"
)

type EventType int

const (
	EventPlayerRegistered EventType = iota + 1
	EventPlayerEntered
	EventMonsterKilled
	EventNextFloor
	EventPreviousFloor
	EventBossFloor
	EventBossKilled
	EventPlayerLeft
	EventCannotContinue
	EventHealthRestored
	EventDamageReceived
)

func (et EventType) String() string {
	names := map[EventType]string{
		EventPlayerRegistered: "PlayerRegistered",
		EventPlayerEntered:    "PlayerEntered",
		EventMonsterKilled:    "MonsterKilled",
		EventNextFloor:        "NextFloor",
		EventPreviousFloor:    "PreviousFloor",
		EventBossFloor:        "BossFloor",
		EventBossKilled:       "BossKilled",
		EventPlayerLeft:       "PlayerLeft",
		EventCannotContinue:   "CannotContinue",
		EventHealthRestored:   "HealthRestored",
		EventDamageReceived:   "DamageReceived",
	}
	if name, ok := names[et]; ok {
		return name
	}
	return "Unknown"
}

func (et EventType) IsValid() bool {
	return et >= EventPlayerRegistered && et <= EventDamageReceived
}

type OutEventType int

const (
	OutEventDisqualified OutEventType = iota + 31
	OutEventDead
	OutEventImpossibleMove
)

type FinalState string

const (
	FinalSuccess FinalState = "SUCCESS"
	FinalFail    FinalState = "FAIL"
	FinalDisqual FinalState = "DISQUAL"
)

const (
	MaxHealth   = 100
	StartHealth = 100
	StartFloor  = 0
)
