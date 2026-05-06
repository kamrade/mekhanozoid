package game

type ActionType string

const (
	ActionTypeStartGame ActionType = "start_game"
	ActionTypeEndTurn   ActionType = "end_turn"
	ActionTypePlayCard  ActionType = "play_card"
	ActionTypeAttack    ActionType = "attack"
)

type TargetType string

const (
	TargetTypeNone   TargetType = "none"
	TargetTypePlayer TargetType = "player"
	TargetTypeBoss   TargetType = "boss"
	TargetTypeMinion TargetType = "minion"
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
	Target   Target
}
