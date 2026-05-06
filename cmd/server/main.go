package main

import (
	"fmt"
	"mekhanozid/internal/game"
)

func main() {
	g := game.NewGame()

	fmt.Println("Mechazod Card Game backend")
	fmt.Printf("Game status: %s\n", g.Status())
}
