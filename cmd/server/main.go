package main

import (
	"fmt"

	"mekhanozid/internal/game"
)

func main() {
	g := game.NewGame()

	fmt.Println("Mechazod Card Game backend")
	fmt.Printf("Game status: %s\n", g.Status)
	fmt.Printf("Players: %d\n", len(g.Players))
	fmt.Printf("Boss: %s\n", g.Boss.Name)
}
