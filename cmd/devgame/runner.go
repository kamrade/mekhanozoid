package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mekhanozid/internal/game"
)

// запускает простой CLI-цикл:
// показывает состояние игры,
// ждёт команды пользователя,
// обрабатывает её,
// потом снова ждёт команду.
func (r *runner) run() {
	fmt.Println("devgame started")
	r.renderState()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n:> ")
		// scanner.Scan() - начинает ожидание ввода.
		// и возвращает false если что-то пошло не так или пользователь нажал Control+C или что-то такое.
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("\ninput error:", err)
			} else {
				fmt.Println("\nexit")
			}
			return
		}

		if r.handleCommand(strings.TrimSpace(scanner.Text())) {
			return
		}
	}
}

// Обрабатываем введеную команду
func (r *runner) handleCommand(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		r.printHelp()
		return false
	}

	command := strings.ToLower(fields[0])
	switch command {
	case "quit", "exit", "q":
		fmt.Println("exit")
		return true
	case "r", "restart":
		r.g = newDevGame()
		fmt.Println("game restarted")
		r.renderState()
	case "?", "help":
		r.printHelp()
	case "s", "state":
		r.renderState()
	case "h", "hand":
		r.renderHand()
	case "history":
		r.renderAllEvents()
	case "e", "end":
		r.cmdEnd(fields)
	case "p", "play":
		r.cmdPlay(fields)
	case "a", "attack":
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
	// Guards
	if len(fields) < 2 || len(fields) > 3 {
		fmt.Println("error: usage: play <handIndex> [targetID]")
		return
	}

	// Получить номер карты
	handIndex, err := strconv.Atoi(fields[1])
	if err != nil {
		fmt.Printf("error: invalid hand index %q\n", fields[1])
		return
	}

	if r.g == nil {
		fmt.Println("error: game is nil")
		return
	}

	// Получить активного игрока
	active := r.activePlayer()
	if handIndex < 1 || handIndex > len(active.Hand) {
		fmt.Printf("error: hand index out of range: %d (valid: 1..%d)\n", handIndex, len(active.Hand))
		return
	}

	// Получить карту
	// Получить targetId, если есть
	card := active.Hand[handIndex-1]
	targetID := ""
	if len(fields) == 3 {
		targetID = fields[2]
	}

	// В итоге здесь мы имеем
	// 1. активного игрока
	// 2. какую карту разыгрываем
	// 3. цель карты (опционально)
	// имея всю эту информацию - пытаемся применить карту
	r.applyAndReport(game.Action{
		Type:     game.ActionPlayCard,
		PlayerID: r.g.ActivePlayerID,
		CardID:   card.ID,
		TargetID: targetID,
	})
}

func (r *runner) cmdAttack(fields []string) {
	// Guards
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
		if errors.Is(err, game.ErrBoardFull) {
			fmt.Println("error: cannot summon more minions. Board limit is 7.")
			return
		}
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("ok: %d event(s)\n", len(events))
	for i := range events {
		e := events[i]
		fmt.Printf("  event: turn=%d type=%s msg=%q amount=%d player=%s\n", e.Turn, e.Type, e.Message, e.Amount, e.PlayerID)
	}

	if r.isGameOver() {
		r.renderAllEvents()
		fmt.Printf("GAME OVER. %s won.\n", r.winnerName())
		fmt.Println("Press r to restart q to exit")
		return
	}

	r.renderState()
}
