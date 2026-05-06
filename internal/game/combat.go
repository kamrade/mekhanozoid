package game

func DealDamage(g *Game, amount int) GameEvent {
	if g == nil {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage: game is nil",
		}
	}

	if amount < 0 {
		amount = 0
	}

	g.Boss.Health -= amount
	if g.Boss.Health < 0 {
		g.Boss.Health = 0
	}

	return GameEvent{
		Type:    EventDamage,
		Target:  Target{Type: TargetTypeBoss, BossID: g.Boss.ID},
		Amount:  amount,
		Message: "damage dealt",
		Turn:    g.Turn,
	}
}
