package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mekhanozid/internal/game"
)

const recentEventsLimit = 8

func main() {
	g := game.NewGame(
		"dev_game_1",
		game.PlayerConfig{ID: "player_1", Name: "Player 1"},
		game.PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	fmt.Println("devgame started")
	printState(g)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			fmt.Println("\nexit")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) == 0 {
			printHelp()
			continue
		}

		command := strings.ToLower(fields[0])

		switch command {
		case "quit", "exit":
			fmt.Println("exit")
			return
		case "help":
			printHelp()
		case "state":
			printState(g)
		case "hand":
			printHand(g)
		case "end":
			if len(fields) == 2 && strings.ToLower(fields[1]) != "turn" {
				fmt.Println("error: usage: end OR end turn")
				continue
			}
			applyAndReport(g, game.Action{
				Type:     game.ActionEndTurn,
				PlayerID: g.ActivePlayerID,
			})
		case "play":
			handlePlayCommand(g, fields)
		case "attack":
			handleAttackCommand(g, fields)
		default:
			fmt.Printf("unknown command: %q\n", line)
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("available actions:")
	fmt.Println("  hand")
	fmt.Println("  play <handIndex> [targetID]")
	fmt.Println("  attack <minionIndex> boss")
	fmt.Println("  end | end turn")
	fmt.Println("  state")
	fmt.Println("  quit | exit")
}

func printState(g *game.Game) {
	if g == nil {
		fmt.Println("game is nil")
		return
	}

	active := findActivePlayer(g)

	fmt.Println("\n=== Game State ===")
	fmt.Printf("status: %s\n", g.Status)
	fmt.Printf("turn: %d\n", g.Turn)
	fmt.Printf("boss: %s hp=%d/%d atk=%d\n", g.Boss.Name, g.Boss.Health, g.Boss.MaxHealth, g.Boss.Attack)

	if len(g.Players) >= 2 {
		fmt.Printf("player 1: %s hp=%d/%d\n", g.Players[0].Name, g.Players[0].Health, g.Players[0].MaxHealth)
		fmt.Printf("player 2: %s hp=%d/%d\n", g.Players[1].Name, g.Players[1].Health, g.Players[1].MaxHealth)
	}

	fmt.Printf("active player: %s (%s)\n", active.Name, active.ID)
	fmt.Printf("active mana: %d/%d\n", active.Mana, active.MaxMana)

	printHand(g)

	fmt.Println("active board:")
	if len(active.Board) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, minion := range active.Board {
			fmt.Printf("  %d. %s (%d/%d) atk=%d canAttack=%t id=%s\n", i+1, minion.Name, minion.Health, minion.MaxHealth, minion.Attack, minion.CanAttack, minion.ID)
		}
	}

	printHelp()
	printRecentEvents(g)
}

func printHand(g *game.Game) {
	if g == nil {
		fmt.Println("active hand:\n  (game is nil)")
		return
	}

	active := findActivePlayer(g)
	fmt.Println("active hand:")
	if len(active.Hand) == 0 {
		fmt.Println("  (empty)")
		return
	}

	for i, card := range active.Hand {
		def, ok := game.CardRegistry[card.DefinitionID]
		if !ok {
			fmt.Printf("  %d. %s (unknown) [%s]\n", i+1, card.DefinitionID, card.ID)
			continue
		}

		targetHint := "no target"
		if def.Targeting.Required {
			targetHint = "target required"
		}

		fmt.Printf(
			"  %d. %s cost=%d type=%s %s [%s]\n",
			i+1,
			def.Name,
			def.Cost,
			def.Type,
			targetHint,
			card.ID,
		)
	}
}

func printRecentEvents(g *game.Game) {
	fmt.Println("recent events:")
	if len(g.Events) == 0 {
		fmt.Println("  (none)")
		return
	}

	start := 0
	if len(g.Events) > recentEventsLimit {
		start = len(g.Events) - recentEventsLimit
	}

	for i := start; i < len(g.Events); i++ {
		e := g.Events[i]
		fmt.Printf("  - turn=%d type=%s msg=%q player=%s amount=%d\n", e.Turn, e.Type, e.Message, e.PlayerID, e.Amount)
	}
}

func findActivePlayer(g *game.Game) game.Player {
	for _, p := range g.Players {
		if p.ID == g.ActivePlayerID {
			return p
		}
	}

	if len(g.Players) > 0 {
		return g.Players[0]
	}

	return game.Player{}
}

func applyAndReport(g *game.Game, action game.Action) {
	events, err := game.ApplyAction(g, action)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("ok: %d event(s)\n", len(events))
	for i := range events {
		e := events[i]
		fmt.Printf("  event: turn=%d type=%s msg=%q amount=%d player=%s\n", e.Turn, e.Type, e.Message, e.Amount, e.PlayerID)
	}

	printState(g)
}

func handlePlayCommand(g *game.Game, fields []string) {
	if len(fields) < 2 || len(fields) > 3 {
		fmt.Println("error: usage: play <handIndex> [targetID]")
		return
	}

	handIndex, err := strconv.Atoi(fields[1])
	if err != nil {
		fmt.Printf("error: invalid hand index %q\n", fields[1])
		return
	}

	if g == nil {
		fmt.Println("error: game is nil")
		return
	}

	active := findActivePlayer(g)
	if handIndex < 1 || handIndex > len(active.Hand) {
		fmt.Printf("error: hand index out of range: %d (valid: 1..%d)\n", handIndex, len(active.Hand))
		return
	}

	card := active.Hand[handIndex-1]

	targetID := ""
	if len(fields) == 3 {
		targetID = fields[2]
	}

	applyAndReport(g, game.Action{
		Type:     game.ActionPlayCard,
		PlayerID: g.ActivePlayerID,
		CardID:   card.ID,
		TargetID: targetID,
	})
}

func handleAttackCommand(g *game.Game, fields []string) {
	if len(fields) != 3 {
		fmt.Println("error: usage: attack <minionIndex> boss")
		return
	}

	if strings.ToLower(fields[2]) != "boss" {
		fmt.Println("error: only boss target is supported: attack <minionIndex> boss")
		return
	}

	minionIndex, err := strconv.Atoi(fields[1])
	if err != nil {
		fmt.Printf("error: invalid minion index %q\n", fields[1])
		return
	}

	if g == nil {
		fmt.Println("error: game is nil")
		return
	}

	active := findActivePlayer(g)
	if minionIndex < 1 || minionIndex > len(active.Board) {
		fmt.Printf("error: minion index out of range: %d (valid: 1..%d)\n", minionIndex, len(active.Board))
		return
	}

	minion := active.Board[minionIndex-1]

	applyAndReport(g, game.Action{
		Type:     game.ActionTypeAttack,
		PlayerID: g.ActivePlayerID,
		SourceID: minion.ID,
		TargetID: game.TargetIDBoss,
	})
}
