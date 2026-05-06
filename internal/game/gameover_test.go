package game

import (
	"errors"
	"testing"
)

// TestStrikeDoesNotWinWhenBossSurvives verifies that playing Strike keeps the
// game active when the boss still has health remaining.
// newTestGameWithStrikeInPlayer1Hand is in play_card_test.go
func TestStrikeDoesNotWinWhenBossSurvives(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 1
	player.MaxMana = 1
	g.Boss.Health = 10

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Status != GameStatusActive {
		t.Fatalf("expected game status %q, got %q", GameStatusActive, g.Status)
	}

	if g.Boss.Health != 7 {
		t.Fatalf("expected boss health 7, got %d", g.Boss.Health)
	}

	if hasEventType(events, EventGameWon) {
		t.Fatalf("expected returned events not to contain %q", EventGameWon)
	}

	if hasEventType(g.Events, EventGameWon) {
		t.Fatalf("expected game events not to contain %q", EventGameWon)
	}
}

// TestStrikeWinsWhenBossHealthReachesZero verifies that playing Strike wins the
// game when the boss reaches zero health.
func TestStrikeWinsWhenBossHealthReachesZero(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 1
	player.MaxMana = 1
	g.Boss.Health = 3

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Status != GameStatusWon {
		t.Fatalf("expected game status %q, got %q", GameStatusWon, g.Status)
	}

	if g.Boss.Health != 0 {
		t.Fatalf("expected boss health to be clamped to 0, got %d", g.Boss.Health)
	}

	if !hasEventType(events, EventGameWon) {
		t.Fatalf("expected returned events to contain %q", EventGameWon)
	}

	if !hasEventType(g.Events, EventGameWon) {
		t.Fatalf("expected game events to contain %q", EventGameWon)
	}
}

// TestStrikeWinsWhenBossHealthGoesBelowZero verifies that boss health is clamped
// to zero even when damage exceeds the boss's remaining health.
func TestStrikeWinsWhenBossHealthGoesBelowZero(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 1
	player.MaxMana = 1
	g.Boss.Health = 1

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Status != GameStatusWon {
		t.Fatalf("expected game status %q, got %q", GameStatusWon, g.Status)
	}

	if g.Boss.Health != 0 {
		t.Fatalf("expected boss health to be clamped to 0, got %d", g.Boss.Health)
	}
}

// TestCheckGameOverDoesNotDuplicateGameWonEvent verifies that calling
// CheckGameOver after the game is already won does not append another win event.
func TestCheckGameOverDoesNotDuplicateGameWonEvent(t *testing.T) {
	g := newTestGame()
	g.Boss.Health = 0

	events := CheckGameOver(g)
	if len(events) != 1 {
		t.Fatalf("expected 1 game over event, got %d", len(events))
	}

	if events[0].Type != EventGameWon {
		t.Fatalf("expected event type %q, got %q", EventGameWon, events[0].Type)
	}

	eventCountAfterFirstCheck := len(g.Events)

	events = CheckGameOver(g)
	if len(events) != 0 {
		t.Fatalf("expected no duplicate game over events, got %d", len(events))
	}

	if len(g.Events) != eventCountAfterFirstCheck {
		t.Fatalf("expected event count to stay %d, got %d", eventCountAfterFirstCheck, len(g.Events))
	}
}

// TestCannotPlayCardAfterGameWon verifies that ApplyAction rejects PlayCard
// actions once the game has already been won.
func TestCannotPlayCardAfterGameWon(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	g.Status = GameStatusWon
	player.Mana = 1
	player.MaxMana = 1

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive, got %v", err)
	}
}

// TestCannotEndTurnAfterGameWon verifies that ApplyAction rejects EndTurn
// actions once the game has already been won.
func TestCannotEndTurnAfterGameWon(t *testing.T) {
	g := newTestGame()

	g.Status = GameStatusWon

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: g.Players[0].ID,
	})

	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive, got %v", err)
	}
}
