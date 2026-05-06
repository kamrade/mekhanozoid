// This file defines player actions and action targets.
// Actions describe player intent and are processed by the game engine through ApplyAction.

package game

type ActionType string

const (
	ActionTypeStartGame ActionType = "start_game"
	ActionTypeEndTurn   ActionType = "end_turn"
	ActionTypePlayCard  ActionType = "play_card"
	ActionTypeAttack    ActionType = "attack"

	// Aliases
	ActionEndTurn  = ActionTypeEndTurn
	ActionPlayCard = ActionTypePlayCard
)

type TargetType string

const (
	TargetTypeNone   TargetType = "none"
	TargetTypePlayer TargetType = "player"
	TargetTypeBoss   TargetType = "boss"
	TargetTypeMinion TargetType = "minion"
)

const (
	TargetIDHero0 = "hero:0"
	TargetIDHero1 = "hero:1"
	TargetIDBoss  = "boss"
)

type Target struct {
	Type     TargetType
	PlayerID PlayerID
	BossID   BossID
	MinionID MinionID
}

type Action struct {
	Type     ActionType
	PlayerID PlayerID
	CardID   CardInstanceID
	SourceID MinionID

	// Почему TargetID отдельно от Target
	// TargetID — входной формат действия: "hero:0", "hero:1", "boss"
	// Target   — нормализованная доменная структура для событий и будущей логики
	TargetID string
	Target   Target
}
