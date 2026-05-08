package game

import "strconv"

// CleanupDeadMinions removes all minions with Health <= 0 from both boards.
// It returns one minion_died event per removed minion.
func CleanupDeadMinions(g *Game) []GameEvent {
	if g == nil {
		return nil
	}

	events := []GameEvent{}

	for playerIndex := range g.Players {
		player := &g.Players[playerIndex]
		living := make([]Minion, 0, len(player.Board))

		for minionIndex := range player.Board {
			minion := player.Board[minionIndex]
			if minion.Health > 0 {
				living = append(living, minion)
				continue
			}

			events = append(events, GameEvent{
				Type: EventMinionDied,
				Target: Target{
					ID:          minionTargetID(playerIndex, minion.ID),
					Type:        TargetTypeMinion,
					Kind:        TargetKindMinion,
					MinionID:    minion.ID,
					OwnerID:     minion.OwnerID,
					DisplayName: minion.Name,
				},
				Message: "minion died",
				Turn:    g.Turn,
			})
		}

		player.Board = living
	}

	return events
}

func minionTargetID(ownerIndex int, minionID MinionID) string {
	return "minion:" + strconv.Itoa(ownerIndex) + ":" + string(minionID)
}
