package main

import (
	"log"
	"mekhanozid/internal/web"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	h := web.NewHandlers()
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/game", h.Game)
	log.Println("server started: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	// g := game.NewGame(
	// 	"game_1",
	// 	game.PlayerConfig{
	// 		ID:   "player_1",
	// 		Name: "Player 1",
	// 	},
	// 	game.PlayerConfig{
	// 		ID:   "player_2",
	// 		Name: "Player 2",
	// 	},
	// 	42,
	// )

	// fmt.Println("Mechazod Card Game backend")
	// fmt.Printf("Game status: %s\n", g.Status)
	// fmt.Printf("Players: %d\n", len(g.Players))
	// fmt.Printf("Boss: %s\n", g.Boss.Name)
	// fmt.Printf("Active player: %s\n", g.ActivePlayerID)
}
