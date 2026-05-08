package game

import (
	"errors"
	"testing"
)

func TestGameCanBePlayedUntilWin(t *testing.T) {
	g := NewGame(
		"scenario_game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	// Arrange a deterministic short scenario: both players can play Strike.
	g.Boss.Health = 6
	g.Players[0].Health = 30
	g.Players[1].Health = 30
	g.Players[0].Mana = 10
	g.Players[0].MaxMana = 10
	g.Players[1].Mana = 10
	g.Players[1].MaxMana = 10
	g.Players[0].Hand = []CardInstance{
		{ID: CardInstanceID("p1_strike_scenario"), DefinitionID: CardID("strike"), OwnerID: g.Players[0].ID},
	}
	g.Players[1].Hand = []CardInstance{
		{ID: CardInstanceID("p2_strike_scenario"), DefinitionID: CardID("strike"), OwnerID: g.Players[1].ID},
	}

	player1ID := g.Players[0].ID
	player2ID := g.Players[1].ID

	if g.ActivePlayerID != player1ID {
		t.Fatalf("expected player 1 to start, got %q", g.ActivePlayerID)
	}

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player1ID,
		CardID:   CardInstanceID("p1_strike_scenario"),
	})
	if err != nil {
		t.Fatalf("expected player 1 strike to succeed, got %v", err)
	}

	if !hasEventType(events, EventCardPlayed) {
		t.Fatalf("expected first action events to contain %q", EventCardPlayed)
	}

	events, err = ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player1ID,
	})
	if err != nil {
		t.Fatalf("expected end turn to succeed, got %v", err)
	}

	if g.ActivePlayerID != player2ID {
		t.Fatalf("expected active player to switch to %q, got %q", player2ID, g.ActivePlayerID)
	}

	if !hasEventType(events, EventBossAbility) {
		t.Fatalf("expected end-turn flow to contain %q", EventBossAbility)
	}

	events, err = ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player2ID,
		CardID:   CardInstanceID("p2_strike_scenario"),
	})
	if err != nil {
		t.Fatalf("expected player 2 strike to succeed, got %v", err)
	}

	if !hasEventType(events, EventCardPlayed) {
		t.Fatalf("expected second strike events to contain %q", EventCardPlayed)
	}

	if g.Boss.Health != 0 {
		t.Fatalf("expected boss health to be clamped to 0, got %d", g.Boss.Health)
	}

	if !hasEventType(g.Events, EventBossAbility) {
		t.Fatalf("expected game events to contain at least one %q", EventBossAbility)
	}

	if !hasEventType(g.Events, EventCardPlayed) {
		t.Fatalf("expected game events to contain %q", EventCardPlayed)
	}

	if g.Status != GameStatusWon {
		t.Fatalf("expected final game status %q, got %q", GameStatusWon, g.Status)
	}
}

func TestGameCanBeLost(t *testing.T) {
	g := NewGame(
		"scenario_game_loss_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		4, // deterministic: at turn 2 boss selects Zap Heroes
	)

	player1ID := g.Players[0].ID
	player2ID := g.Players[1].ID

	// Arrange lethal setup for boss turn-start ability.
	g.Players[0].Health = 2
	g.Players[1].Health = 10

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player1ID,
	})
	if err != nil {
		t.Fatalf("expected end turn to succeed, got %v", err)
	}

	if g.Players[0].Health != 0 {
		t.Fatalf("expected player 1 health to be clamped to 0, got %d", g.Players[0].Health)
	}

	if g.Status != GameStatusLost {
		t.Fatalf("expected final game status %q, got %q", GameStatusLost, g.Status)
	}

	if !hasEventType(events, EventGameLost) {
		t.Fatalf("expected returned events to contain %q", EventGameLost)
	}

	if !hasEventType(g.Events, EventGameLost) {
		t.Fatalf("expected game events to contain %q", EventGameLost)
	}

	_, err = ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: player2ID,
	})
	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive for EndTurn after loss, got %v", err)
	}

	_, err = ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player2ID,
		CardID:   CardInstanceID("any_card"),
	})
	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive for PlayCard after loss, got %v", err)
	}

	_, err = ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player2ID,
		SourceID: MinionID("any_minion"),
		TargetID: TargetIDBoss,
	})
	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive for Attack after loss, got %v", err)
	}
}
