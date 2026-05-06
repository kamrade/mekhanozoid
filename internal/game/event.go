// This file defines game events.
// Events describe important changes in game state and can later be used by UI, logs, or replays.

package game

type EventType string

const (
	EventTypeGameCreated EventType = "game_created"
	EventTypeGameStarted EventType = "game_started"
	EventTypeTurnStarted EventType = "turn_started"
	EventTypeTurnEnded   EventType = "turn_ended"
	EventTypeCardPlayed  EventType = "card_played"
	EventTypeDamageDealt EventType = "damage_dealt"
	EventTypeHeal        EventType = "heal"
	EventTypeGameWon     EventType = "game_won"
	EventTypeGameLost    EventType = "game_lost"
	EventTypeCardDrawn   EventType = "card_drawn"

	// Aliases
	EventCardPlayed = EventTypeCardPlayed
	EventDamage     = EventTypeDamageDealt
	EventHeal       = EventTypeHeal
	EventGameWon    = EventTypeGameWon
)

type GameEvent struct {
	Type     EventType
	PlayerID PlayerID
	CardID   CardInstanceID
	SourceID MinionID
	Target   Target
	Amount   int
	Message  string
	Turn     int
}
