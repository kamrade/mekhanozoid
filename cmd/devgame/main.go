package main

import "mekhanozid/internal/game"

type runner struct {
	g *game.Game
}

func main() {
	r := &runner{g: newDevGame()}
	r.run()
}
