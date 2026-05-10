package main

import (
	"fmt"

	"mekhanozid/internal/game"
)

func (r *runner) renderState() {
	if r.g == nil {
		fmt.Println("game is nil")
		return
	}

	active := r.activePlayer()

	fmt.Println()
	fmt.Println("\n=== === === === === Game State === === === === ===")
	fmt.Printf("status              : %s\n", r.g.Status)
	if r.g.Status == game.GameStatusWon {
		fmt.Println("GAME OVER")
		fmt.Println("winner: players")
	}
	if r.g.Status == game.GameStatusLost {
		fmt.Println("GAME OVER")
		fmt.Println("winner: boss")
	}
	fmt.Printf("turn                : %d\n", r.g.Turn)
	fmt.Printf(
		"boss                : %s hp=%s%d/%d%s atk=%d\n",
		r.g.Boss.Name,
		ColorGreen,
		r.g.Boss.Health,
		r.g.Boss.MaxHealth,
		ColorReset,
		r.g.Boss.Attack,
	)

	if len(r.g.Players) >= 2 {
		hpColor1 := ColorGreen
		if r.g.Players[0].Health < 10 {
			hpColor1 = ColorRed
		}

		fmt.Printf(
			"player 1            : %s hp=%s%d/%d%s deck=%d\n",
			r.g.Players[0].Name,
			hpColor1,
			r.g.Players[0].Health,
			r.g.Players[0].MaxHealth,
			ColorReset,
			len(r.g.Players[0].Deck),
		)

		hpColor2 := ColorGreen
		if r.g.Players[1].Health < 10 {
			hpColor2 = ColorRed
		}
		fmt.Printf(
			"player 2            : %s hp=%s%d/%d%s deck=%d\n",
			r.g.Players[1].Name,
			hpColor2,
			r.g.Players[1].Health,
			r.g.Players[1].MaxHealth,
			ColorReset,
			len(r.g.Players[1].Deck),
		)
	}

	fmt.Printf(ColorYellow+"active player       : %s (%s)\n"+ColorReset, active.Name, active.ID)
	fmt.Printf("active mana         : %s%d/%d%s\n", ColorBlue, active.Mana, active.MaxMana, ColorReset)

	r.renderHand()
	r.renderBoard(active)
	r.renderOtherBoard(active.ID)
	fmt.Println("---")
	// r.printHelp()
	r.renderRecentEvents()
}

func (r *runner) renderHand() {
	if r.g == nil {
		fmt.Print("active hand         :\n  (game is nil)\n")
		return
	}

	active := r.activePlayer()
	fmt.Println("active hand         :")
	if len(active.Hand) == 0 {
		fmt.Println("                      (empty)")
		return
	}

	for i, card := range active.Hand {
		def, ok := game.CardRegistry[card.DefinitionID]
		if !ok {
			fmt.Printf("                      %d. %s (unknown) [%s]\n", i+1, card.DefinitionID, card.ID)
			continue
		}

		targetHint := "no target"
		if def.Targeting.Required {
			targetHint = "target required"
		}

		fmt.Printf("                      %d. %s cost=%d type=%s %s [%s]\n", i+1, def.Name, def.Cost, def.Type, targetHint, card.ID)
	}
}

func (r *runner) renderBoard(active game.Player) {
	fmt.Println("active board:       ")
	if len(active.Board) == 0 {
		fmt.Println("                      (empty)")
		return
	}

	for i, minion := range active.Board {
		attackedThisTurn := !minion.CanAttack

		canAttackColor := ColorGreen
		if !minion.CanAttack {
			canAttackColor = ColorRed
		}

		fmt.Printf(
			"                      %d. %s (%d/%d) atk=%d canAttack=%s%t%s attackedThisTurn=%t id=%s\n",
			i+1,
			minion.Name,
			minion.Health,
			minion.MaxHealth,
			minion.Attack,
			canAttackColor,
			minion.CanAttack,
			ColorReset,
			attackedThisTurn,
			minion.ID,
		)
	}
}

func (r *runner) renderOtherBoard(activePlayerID game.PlayerID) {
	if r.g == nil {
		fmt.Println("other board:")
		fmt.Println("  (game is nil)")
		return
	}

	fmt.Println("other board:")
	for i := range r.g.Players {
		if r.g.Players[i].ID == activePlayerID {
			continue
		}

		board := r.g.Players[i].Board
		if len(board) == 0 {
			fmt.Println("                      (empty)")
			return
		}

		for j, minion := range board {
			attackedThisTurn := !minion.CanAttack
			fmt.Printf(
				"                      %d. %s (%d/%d) atk=%d canAttack=%t attackedThisTurn=%t id=%s\n",
				j+1,
				minion.Name,
				minion.Health,
				minion.MaxHealth,
				minion.Attack,
				minion.CanAttack,
				attackedThisTurn,
				minion.ID,
			)
		}
		return
	}

	fmt.Println("  (no other player)")
}

func (r *runner) renderRecentEvents() {
	fmt.Println("recent events:      ")
	if r.g == nil || len(r.g.Events) == 0 {
		fmt.Println("                      (none)")
		return
	}

	start := 0
	if len(r.g.Events) > recentEventsLimit {
		start = len(r.g.Events) - recentEventsLimit
	}

	for i := start; i < len(r.g.Events); i++ {
		e := r.g.Events[i]
		fmt.Printf("                    - turn=%d type=%s %s%q%s player=%s amount=%d\n", e.Turn, e.Type, ColorBlue, e.Message, ColorReset, e.PlayerID, e.Amount)
	}
}

func (r *runner) renderAllEvents() {
	fmt.Println("event history:")
	if r.g == nil || len(r.g.Events) == 0 {
		fmt.Println("  (none)")
		return
	}

	for i := range r.g.Events {
		e := r.g.Events[i]
		fmt.Printf("  %d. turn=%d type=%s msg=%q player=%s amount=%d\n", i+1, e.Turn, e.Type, e.Message, e.PlayerID, e.Amount)
	}
}

func (r *runner) printHelp() {
	fmt.Println("available actions:")
	fmt.Println("                      hand (h)")
	fmt.Println("                      p | play <handIndex> [targetID]")
	fmt.Println("                      a | attack <minionIndex> boss")
	fmt.Println("                      e | end")
	fmt.Println("                      s | state")
	fmt.Println("                      ? | help")
	fmt.Println("                      history")
	fmt.Println("                      q | quit | exit")
	fmt.Println("                      r | restart")
}
