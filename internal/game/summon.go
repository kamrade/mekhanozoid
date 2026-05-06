package game

import "fmt"

func canSummonMinion(player *Player) error {
	if player == nil {
		return ErrInvalidPlayerIndex
	}

	if len(player.Board) >= MaxBoardSize {
		return ErrBoardFull
	}

	return nil
}

// SummonMinion creates a board minion from a minion card definition and adds it to the player's board.
func SummonMinion(g *Game, playerIndex int, card CardDefinition, cardInstance CardInstance) (GameEvent, error) {
	if g == nil {
		return GameEvent{}, ErrNilGame
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return GameEvent{}, ErrInvalidPlayerIndex
	}

	player := &g.Players[playerIndex]

	if err := canSummonMinion(player); err != nil {
		return GameEvent{}, err
	}

	minion := Minion{
		ID:           MinionID(fmt.Sprintf("%s_minion_%d", player.ID, len(player.Board)+1)),
		DefinitionID: card.ID,
		OwnerID:      player.ID,
		Name:         card.Name,
		Attack:       card.Attack,
		Health:       card.Health,
		MaxHealth:    card.Health,
		CanAttack:    false,
		Exhausted:    true,
	}

	player.Board = append(player.Board, minion)

	return GameEvent{
		Type:     EventMinionSummoned,
		PlayerID: player.ID,
		CardID:   cardInstance.ID,
		SourceID: minion.ID,
		Target: Target{
			Type:     TargetTypePlayer,
			Kind:     TargetKindHero,
			PlayerID: player.ID,
			OwnerID:  player.ID,
		},
		Message: "minion summoned",
		Turn:    g.Turn,
	}, nil
}
