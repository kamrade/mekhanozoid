package game

// CheckGameOver checks whether the game has reached a terminal state.
// Loss has precedence over win: if any player is dead, the game is lost.
func CheckGameOver(g *Game) []GameEvent {
	if g == nil {
		return nil
	}

	if g.Status != GameStatusActive {
		return nil
	}

	for i := range g.Players {
		if g.Players[i].Health < 0 {
			g.Players[i].Health = 0
		}

		if g.Players[i].Health == 0 {
			g.Status = GameStatusLost

			event := GameEvent{
				Type:     EventGameLost,
				PlayerID: g.Players[i].ID,
				Target: Target{
					Type:        TargetTypePlayer,
					Kind:        TargetKindHero,
					PlayerID:    g.Players[i].ID,
					OwnerID:     g.Players[i].ID,
					DisplayName: g.Players[i].Name,
				},
				Message: "game lost",
				Turn:    g.Turn,
			}

			g.Events = append(g.Events, event)
			return []GameEvent{event}
		}
	}

	if g.Boss.Health > 0 {
		return nil
	}

	g.Status = GameStatusWon

	event := GameEvent{
		Type:    EventGameWon,
		Target:  Target{Type: TargetTypeBoss, BossID: g.Boss.ID},
		Message: "game won",
		Turn:    g.Turn,
	}

	g.Events = append(g.Events, event)

	return []GameEvent{event}
}
