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
const devBossHealth = 30

type runner struct {
	g *game.Game
}

func main() {
	r := &runner{g: newDevGame()}
	r.run()
}

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

func (r *runner) run() {
	fmt.Println("devgame started")
	r.renderState()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			fmt.Println("\nexit")
			return
		}

		if r.handleCommand(strings.TrimSpace(scanner.Text())) {
			return
		}
	}
}

func (r *runner) handleCommand(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		r.printHelp()
		return false
	}

	command := strings.ToLower(fields[0])
	switch command {
	case "quit", "exit":
		fmt.Println("exit")
		return true
	case "help":
		r.printHelp()
	case "state":
		r.renderState()
	case "hand":
		r.renderHand()
	case "end":
		r.cmdEnd(fields)
	case "play":
		r.cmdPlay(fields)
	case "attack":
		r.cmdAttack(fields)
	default:
		fmt.Printf("unknown command: %q\n", line)
		r.printHelp()
	}

	return false
}

func (r *runner) cmdEnd(fields []string) {
	if len(fields) == 2 && strings.ToLower(fields[1]) != "turn" {
		fmt.Println("error: usage: end OR end turn")
		return
	}

	r.applyAndReport(game.Action{
		Type:     game.ActionEndTurn,
		PlayerID: r.g.ActivePlayerID,
	})
}

func (r *runner) cmdPlay(fields []string) {
	if len(fields) < 2 || len(fields) > 3 {
		fmt.Println("error: usage: play <handIndex> [targetID]")
		return
	}

	handIndex, err := strconv.Atoi(fields[1])
	if err != nil {
		fmt.Printf("error: invalid hand index %q\n", fields[1])
		return
	}

	if r.g == nil {
		fmt.Println("error: game is nil")
		return
	}

	active := r.activePlayer()
	if handIndex < 1 || handIndex > len(active.Hand) {
		fmt.Printf("error: hand index out of range: %d (valid: 1..%d)\n", handIndex, len(active.Hand))
		return
	}

	card := active.Hand[handIndex-1]
	targetID := ""
	if len(fields) == 3 {
		targetID = fields[2]
	}

	r.applyAndReport(game.Action{
		Type:     game.ActionPlayCard,
		PlayerID: r.g.ActivePlayerID,
		CardID:   card.ID,
		TargetID: targetID,
	})
}

func (r *runner) cmdAttack(fields []string) {
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

	if r.g == nil {
		fmt.Println("error: game is nil")
		return
	}

	active := r.activePlayer()
	if minionIndex < 1 || minionIndex > len(active.Board) {
		fmt.Printf("error: minion index out of range: %d (valid: 1..%d)\n", minionIndex, len(active.Board))
		return
	}

	minion := active.Board[minionIndex-1]
	r.applyAndReport(game.Action{
		Type:     game.ActionTypeAttack,
		PlayerID: r.g.ActivePlayerID,
		SourceID: minion.ID,
		TargetID: game.TargetIDBoss,
	})
}

func (r *runner) applyAndReport(action game.Action) {
	events, err := game.ApplyAction(r.g, action)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("ok: %d event(s)\n", len(events))
	for i := range events {
		e := events[i]
		fmt.Printf("  event: turn=%d type=%s msg=%q amount=%d player=%s\n", e.Turn, e.Type, e.Message, e.Amount, e.PlayerID)
	}

	r.renderState()
}

func (r *runner) renderState() {
	if r.g == nil {
		fmt.Println("game is nil")
		return
	}

	active := r.activePlayer()

	fmt.Println("\n")
	fmt.Println("\n=== Game State ===")
	fmt.Printf("status: %s\n", r.g.Status)
	if r.g.Status == game.GameStatusWon {
		fmt.Println("GAME OVER: PLAYERS WON")
	}
	if r.g.Status == game.GameStatusLost {
		fmt.Println("GAME OVER: PLAYERS LOST")
	}
	fmt.Printf("turn: %d\n", r.g.Turn)
	fmt.Printf("boss: %s hp=%d/%d atk=%d\n", r.g.Boss.Name, r.g.Boss.Health, r.g.Boss.MaxHealth, r.g.Boss.Attack)

	if len(r.g.Players) >= 2 {
		fmt.Printf("player 1: %s hp=%d/%d\n", r.g.Players[0].Name, r.g.Players[0].Health, r.g.Players[0].MaxHealth)
		fmt.Printf("player 2: %s hp=%d/%d\n", r.g.Players[1].Name, r.g.Players[1].Health, r.g.Players[1].MaxHealth)
	}

	fmt.Printf("active player: %s (%s)\n", active.Name, active.ID)
	fmt.Printf("active mana: %d/%d\n", active.Mana, active.MaxMana)

	r.renderHand()
	r.renderBoard(active)
	fmt.Println("---")
	r.printHelp()
	r.renderRecentEvents()
}

func (r *runner) renderHand() {
	if r.g == nil {
		fmt.Println("active hand:\n  (game is nil)")
		return
	}

	active := r.activePlayer()
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

		fmt.Printf("  %d. %s cost=%d type=%s %s [%s]\n", i+1, def.Name, def.Cost, def.Type, targetHint, card.ID)
	}
}

func (r *runner) renderBoard(active game.Player) {
	fmt.Println("active board:")
	if len(active.Board) == 0 {
		fmt.Println("  (empty)")
		return
	}

	for i, minion := range active.Board {
		fmt.Printf("  %d. %s (%d/%d) atk=%d canAttack=%t id=%s\n", i+1, minion.Name, minion.Health, minion.MaxHealth, minion.Attack, minion.CanAttack, minion.ID)
	}
}

func (r *runner) renderRecentEvents() {
	fmt.Println("recent events:")
	if r.g == nil || len(r.g.Events) == 0 {
		fmt.Println("  (none)")
		return
	}

	start := 0
	if len(r.g.Events) > recentEventsLimit {
		start = len(r.g.Events) - recentEventsLimit
	}

	for i := start; i < len(r.g.Events); i++ {
		e := r.g.Events[i]
		fmt.Printf("  - turn=%d type=%s msg=%q player=%s amount=%d\n", e.Turn, e.Type, e.Message, e.PlayerID, e.Amount)
	}

}

func (r *runner) printHelp() {
	fmt.Println("available actions:")
	fmt.Println("  hand")
	fmt.Println("  play <handIndex> [targetID]")
	fmt.Println("  attack <minionIndex> boss")
	fmt.Println("  end | end turn")
	fmt.Println("  state")
	fmt.Println("  quit | exit")
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

func prepareDevGame(g *game.Game) {
	if g == nil {
		return
	}

	if devBossHealth > 0 && devBossHealth < g.Boss.MaxHealth {
		g.Boss.Health = devBossHealth
		g.Boss.MaxHealth = devBossHealth
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
