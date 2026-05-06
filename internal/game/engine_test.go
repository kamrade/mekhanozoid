package game

import (
	"errors"
	"testing"
)

// TestApplyActionEndTurnSwitchesFromPlayer1ToPlayer2 verifies that EndTurn
// moves the active turn from the first player to the second player and applies start-of-turn effects.
func TestApplyActionEndTurnSwitchesFromPlayer1ToPlayer2(t *testing.T) {
	g := newTestGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	initialTurn := g.Turn
	initialPlayer2DeckSize := len(g.Players[1].Deck)
	initialPlayer2HandSize := len(g.Players[1].Hand)

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player1.ID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.ActivePlayerID != player2.ID {
		t.Fatalf("expected active player %q, got %q", player2.ID, g.ActivePlayerID)
	}

	if g.Players[0].IsCurrent {
		t.Fatal("expected player 1 to be inactive")
	}

	if !g.Players[1].IsCurrent {
		t.Fatal("expected player 2 to be active")
	}

	if g.Turn != initialTurn+1 {
		t.Fatalf("expected turn %d, got %d", initialTurn+1, g.Turn)
	}

	if g.Players[1].MaxMana != 1 {
		t.Fatalf("expected player 2 max mana 1, got %d", g.Players[1].MaxMana)
	}

	if g.Players[1].Mana != g.Players[1].MaxMana {
		t.Fatalf("expected player 2 mana to equal max mana, got mana=%d max=%d", g.Players[1].Mana, g.Players[1].MaxMana)
	}

	if len(g.Players[1].Deck) != initialPlayer2DeckSize-1 {
		t.Fatalf("expected player 2 deck size %d, got %d", initialPlayer2DeckSize-1, len(g.Players[1].Deck))
	}

	if len(g.Players[1].Hand) != initialPlayer2HandSize+1 {
		t.Fatalf("expected player 2 hand size %d, got %d", initialPlayer2HandSize+1, len(g.Players[1].Hand))
	}

	if !hasEventType(events, EventTypeTurnStarted) {
		t.Fatalf("expected returned events to contain %q", EventTypeTurnStarted)
	}

	if !hasEventType(events, EventTypeCardDrawn) {
		t.Fatalf("expected returned events to contain %q", EventTypeCardDrawn)
	}
}

// TestApplyActionEndTurnSwitchesBackToPlayer1 verifies that consecutive EndTurn
// actions rotate the active player back to the first player.
func TestApplyActionEndTurnSwitchesBackToPlayer1(t *testing.T) {
	g := newTestGame()

	player1ID := g.Players[0].ID
	player2ID := g.Players[1].ID

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player1ID,
	})
	if err != nil {
		t.Fatalf("expected first end turn to succeed, got %v", err)
	}

	_, err = ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player2ID,
	})
	if err != nil {
		t.Fatalf("expected second end turn to succeed, got %v", err)
	}

	if g.ActivePlayerID != player1ID {
		t.Fatalf("expected active player %q, got %q", player1ID, g.ActivePlayerID)
	}

	if !g.Players[0].IsCurrent {
		t.Fatal("expected player 1 to be active")
	}

	if g.Players[1].IsCurrent {
		t.Fatal("expected player 2 to be inactive")
	}
}

// TestApplyActionEndTurnIncreasesTurnEachTime verifies that every successful
// EndTurn action increments the game's turn counter exactly once.
func TestApplyActionEndTurnIncreasesTurnEachTime(t *testing.T) {
	g := newTestGame()

	initialTurn := g.Turn

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: g.Players[0].ID,
	})
	if err != nil {
		t.Fatalf("expected first end turn to succeed, got %v", err)
	}

	if g.Turn != initialTurn+1 {
		t.Fatalf("expected turn %d, got %d", initialTurn+1, g.Turn)
	}

	_, err = ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: g.Players[1].ID,
	})
	if err != nil {
		t.Fatalf("expected second end turn to succeed, got %v", err)
	}

	if g.Turn != initialTurn+2 {
		t.Fatalf("expected turn %d, got %d", initialTurn+2, g.Turn)
	}
}

// TestApplyActionEndTurnAppendsEventsToGameEvents verifies that events produced
// by ApplyAction are also stored in the game's event history.
func TestApplyActionEndTurnAppendsEventsToGameEvents(t *testing.T) {
	g := newTestGame()

	initialEventCount := len(g.Events)

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: g.Players[0].ID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Fatal("expected returned events")
	}

	if len(g.Events) < initialEventCount+len(events) {
		t.Fatalf(
			"expected game events to contain at least %d new events, got initial=%d current=%d returned=%d",
			len(events),
			initialEventCount,
			len(g.Events),
			len(events),
		)
	}

	if !hasEventType(g.Events, EventTypeTurnStarted) {
		t.Fatalf("expected game events to contain %q", EventTypeTurnStarted)
	}

	if !hasEventType(g.Events, EventTypeCardDrawn) {
		t.Fatalf("expected game events to contain %q", EventTypeCardDrawn)
	}
}

// TestApplyActionEndTurnRejectsInactivePlayer verifies that a player who is not
// currently active cannot end another player's turn.
func TestApplyActionEndTurnRejectsInactivePlayer(t *testing.T) {
	g := newTestGame()

	inactivePlayer := g.Players[1]

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: inactivePlayer.ID,
	})

	if !errors.Is(err, ErrNotYourTurn) {
		t.Fatalf("expected ErrNotYourTurn, got %v", err)
	}
}

// TestApplyActionReturnsErrorWhenGameIsNotActive verifies that actions are
// rejected when the game is not in the active status.
func TestApplyActionReturnsErrorWhenGameIsNotActive(t *testing.T) {
	statuses := []GameStatus{
		GameStatusCreated,
		GameStatusWon,
		GameStatusLost,
	}

	for _, status := range statuses {
		g := newTestGame()
		g.Status = status

		_, err := ApplyAction(g, Action{
			Type:     ActionTypeEndTurn,
			PlayerID: g.Players[0].ID,
		})

		if !errors.Is(err, ErrGameNotActive) {
			t.Fatalf("expected ErrGameNotActive for status %q, got %v", status, err)
		}
	}
}

// TestApplyActionReturnsErrorForUnknownAction verifies that unsupported action
// types are rejected instead of being silently ignored.
func TestApplyActionReturnsErrorForUnknownAction(t *testing.T) {
	g := newTestGame()

	_, err := ApplyAction(g, Action{
		Type:     ActionType("unknown_action"),
		PlayerID: g.Players[0].ID,
	})

	if !errors.Is(err, ErrUnknownAction) {
		t.Fatalf("expected ErrUnknownAction, got %v", err)
	}
}

// TestApplyActionReturnsErrorForNilGame verifies that ApplyAction handles a nil
// game pointer safely and returns ErrNilGame.
func TestApplyActionReturnsErrorForNilGame(t *testing.T) {
	_, err := ApplyAction(nil, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: PlayerID("player_1"),
	})

	if !errors.Is(err, ErrNilGame) {
		t.Fatalf("expected ErrNilGame, got %v", err)
	}
}

// TestApplyActionReturnsErrorForUnknownPlayer verifies that actions from a
// player ID that does not exist in the game are rejected.
func TestApplyActionReturnsErrorForUnknownPlayer(t *testing.T) {
	g := newTestGame()

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: PlayerID("missing_player"),
	})

	if !errors.Is(err, ErrInvalidPlayerIndex) {
		t.Fatalf("expected ErrInvalidPlayerIndex, got %v", err)
	}
}

// newTestGame creates a deterministic two-player game for engine tests.
func newTestGame() *Game {
	return NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)
}

// hasEventType reports whether a slice of game events contains at least one
// event with the requested type.
func hasEventType(events []GameEvent, eventType EventType) bool {
	for _, event := range events {
		if event.Type == eventType {
			return true
		}
	}

	return false
}
