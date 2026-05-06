package game

// CheckGameOver checks whether the game has reached a terminal state.
// At this stage, only the win condition is implemented: the boss reaches 0 health.
func CheckGameOver(g *Game) []GameEvent {
	if g == nil {
		return nil
	}

	if g.Status != GameStatusActive {
		return nil
	}

	if g.Boss.Health > 0 {
		return nil
	}

	if g.Boss.Health < 0 {
		g.Boss.Health = 0
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
