package main

import (
	"fmt"

	"mekhanozid/internal/game"
)

func main() {
	g := game.NewGame(
		"game_1",
		game.PlayerConfig{
			ID:   "player_1",
			Name: "Player 1",
		},
		game.PlayerConfig{
			ID:   "player_2",
			Name: "Player 2",
		},
		42,
	)

	fmt.Println("Mechazod Card Game backend")
	fmt.Printf("Game status: %s\n", g.Status)
	fmt.Printf("Players: %d\n", len(g.Players))
	fmt.Printf("Boss: %s\n", g.Boss.Name)
	fmt.Printf("Active player: %s\n", g.ActivePlayerID)
}
