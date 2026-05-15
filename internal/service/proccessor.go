package service

import (
	"dungeon/internal/domain/models"
	"fmt"
	"time"
)

type Processor struct {
	dungeon   *models.Dungeon
	players   map[int]*models.Player
	validator *Validator
	reporter  *Reporter
	output    []string
}

func NewProcessor(dungeon *models.Dungeon) *Processor {
	return &Processor{
		dungeon:   dungeon,
		players:   make(map[int]*models.Player),
		validator: NewValidator(dungeon),
		reporter:  NewReporter(dungeon),
		output:    make([]string, 0),
	}
}

// ProcessEvent takes an Event, retrieves or creates the corresponding Player, validates the event against the player's current state and the dungeon's rules, and if valid, updates the player's state accordingly while generating output messages for each processed event
func (p *Processor) ProcessEvent(event *models.Event) {
	player := p.getOrCreatePlayer(event.PlayerID)

	if err := p.validator.Validate(event, player); err != nil {
		p.addOutput(err.Error())
		return
	}

	p.handleEvent(event, player)
}

// getOrCreatePlayer retrieves an existing Player by ID or creates a new one if it doesn't exist,
// ensuring that all events are associated with a valid Player instance
func (p *Processor) getOrCreatePlayer(id int) *models.Player {
	if player, exists := p.players[id]; exists {
		return player
	}

	player := models.NewPlayer(id)
	p.players[id] = player
	return player
}

// handleEvent updates the player's state based on the type of event, performing actions such as registering,
// entering the dungeon, killing monsters, moving floors, and leaving the dungeon,
// while generating appropriate output messages for each action taken
func (p *Processor) handleEvent(event *models.Event, player *models.Player) {
	t := event.Time

	switch event.Type {
	case models.EventPlayerRegistered:
		player.Register()
		p.addOutput(fmt.Sprintf("[%s] Player [%d] registered", formatTime(t), player.ID))

	case models.EventPlayerEntered:
		player.EnterDungeon(t)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] entered the dungeon", formatTime(t), player.ID))

	case models.EventMonsterKilled:
		player.KillMonster()

		monstersOnFloor := p.dungeon.MonstersOnFloor(player.CurrentFloor)
		if player.MonstersKilled%monstersOnFloor == 0 {
			player.CompleteFloor(t)
		}

		p.addOutput(fmt.Sprintf("[%s] Player [%d] killed the monster", formatTime(t), player.ID))

	case models.EventNextFloor:
		if player.MoveToNextFloor(t, p.dungeon.Floors) {
			p.addOutput(fmt.Sprintf("[%s] Player [%d] went to the next floor", formatTime(t), player.ID))
		}

	case models.EventPreviousFloor:
		if player.MoveToPreviousFloor(t) {
			p.addOutput(fmt.Sprintf("[%s] Player [%d] went to the previous floor", formatTime(t), player.ID))
		}

	case models.EventBossFloor:
		player.EnterBossFloor(t)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] entered the boss's floor", formatTime(t), player.ID))

	case models.EventBossKilled:
		player.KillBoss(t)
		player.CompleteFloor(t)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] killed the boss", formatTime(t), player.ID))

	case models.EventPlayerLeft:
		player.LeaveDungeon(t)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] left the dungeon", formatTime(t), player.ID))

	case models.EventCannotContinue:
		player.Disqualify()
		p.addOutput(fmt.Sprintf("[%s] Player [%d] cannot continue due to [%s]",
			formatTime(t), player.ID, event.ExtraParam))

	case models.EventHealthRestored:
		amount := parseIntParam(event.ExtraParam)
		player.Heal(amount)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] has restored [%s] of health",
			formatTime(t), player.ID, event.ExtraParam))

	case models.EventDamageReceived:
		amount := parseIntParam(event.ExtraParam)
		player.TakeDamage(amount)
		p.addOutput(fmt.Sprintf("[%s] Player [%d] recieved [%s] of damage",
			formatTime(t), player.ID, event.ExtraParam))

		if !player.IsAlive() {
			p.addOutput(fmt.Sprintf("[%s] Player [%d] is dead", formatTime(t), player.ID))
		}
	}
}

func (p *Processor) PrintOutput() {
	for _, msg := range p.output {
		fmt.Println(msg)
	}
}

func (p *Processor) PrintFinalReport() {
	fmt.Println("\nFinal report:")
	for _, player := range p.players {
		report := p.reporter.Generate(player)
		fmt.Println(report)
	}
}

func (p *Processor) addOutput(msg string) {
	p.output = append(p.output, msg)
}

func formatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func parseIntParam(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
