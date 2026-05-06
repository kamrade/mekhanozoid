package game

// Heal restores health to a player and clamps the result to MaxHealth.
func Heal(g *Game, playerIndex int, amount int) GameEvent {
	if g == nil {
		return GameEvent{
			Type:    EventHeal,
			Amount:  0,
			Message: "cannot heal: game is nil",
		}
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return GameEvent{
			Type:    EventHeal,
			Amount:  0,
			Message: "cannot heal: invalid player index",
			Turn:    g.Turn,
		}
	}

	if amount < 0 {
		amount = 0
	}

	player := &g.Players[playerIndex]
	before := player.Health

	player.Health += amount
	if player.Health > player.MaxHealth {
		player.Health = player.MaxHealth
	}

	healedAmount := player.Health - before

	return GameEvent{
		Type:     EventHeal,
		PlayerID: player.ID,
		Target: Target{
			Type:     TargetTypePlayer,
			PlayerID: player.ID,
		},
		Amount:  healedAmount,
		Message: "health restored",
		Turn:    g.Turn,
	}
}
