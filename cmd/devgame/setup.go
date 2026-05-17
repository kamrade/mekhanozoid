package main

import "mekhanozid/internal/game"

func newDevGame() *game.Game {
	g := game.NewGame(
		"dev_game_1",
		game.PlayerConfig{ID: "player_1", Name: "Player 1"},
		game.PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)
	prepareDevGame(g)
	return g
}

func prepareDevGame(g *game.Game) {
	if g == nil {
		return
	}

	if DevBossHealth > 0 && DevBossHealth < g.Boss.MaxHealth {
		g.Boss.Health = DevBossHealth
		g.Boss.MaxHealth = DevBossHealth
	}

	for i := range g.Players {
		if g.Players[i].MaxMana < 1 {
			g.Players[i].MaxMana = 1
		}
		if g.Players[i].Mana < 1 {
			g.Players[i].Mana = 1
		}
	}
}

func (r *runner) activePlayer() game.Player {
	if r.g == nil {
		return game.Player{}
	}

	for _, p := range r.g.Players {
		if p.ID == r.g.ActivePlayerID {
			return p
		}
	}

	if len(r.g.Players) > 0 {
		return r.g.Players[0]
	}

	return game.Player{}
}

func (r *runner) isGameOver() bool {
	if r.g == nil {
		return false
	}

	return r.g.Status == game.GameStatusWon || r.g.Status == game.GameStatusLost
}

func (r *runner) winnerName() string {
	if r.g == nil {
		return "unknown"
	}

	if r.g.Status == game.GameStatusWon {
		return "players"
	}

	if r.g.Status == game.GameStatusLost {
		return "boss"
	}

	return "unknown"
}
