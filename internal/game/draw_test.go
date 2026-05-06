package game

import "testing"

// TestDrawCardMovesTopCardFromDeckToHand verifies that DrawCard takes the top
// card from a player's deck and appends that exact card to the player's hand.
func TestDrawCardMovesTopCardFromDeckToHand(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	playerIndex := 0
	player := &g.Players[playerIndex]

	initialDeckSize := len(player.Deck)
	initialHandSize := len(player.Hand)
	topCard := player.Deck[0]

	events := DrawCard(g, playerIndex)

	if len(player.Deck) != initialDeckSize-1 {
		t.Fatalf("expected deck size %d, got %d", initialDeckSize-1, len(player.Deck))
	}

	if len(player.Hand) != initialHandSize+1 {
		t.Fatalf("expected hand size %d, got %d", initialHandSize+1, len(player.Hand))
	}

	drawnCard := player.Hand[len(player.Hand)-1]

	if drawnCard != topCard {
		t.Fatalf("expected drawn card %+v, got %+v", topCard, drawnCard)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Type != EventTypeCardDrawn {
		t.Fatalf("expected event type %q, got %q", EventTypeCardDrawn, events[0].Type)
	}

	if events[0].PlayerID != player.ID {
		t.Fatalf("expected event player ID %q, got %q", player.ID, events[0].PlayerID)
	}

	if events[0].CardID != topCard.ID {
		t.Fatalf("expected event card ID %q, got %q", topCard.ID, events[0].CardID)
	}
}

// TestDrawCardAppendsEventToGameEvents verifies that a successful card draw
// returns an event and also stores that event in the game's event history.
func TestDrawCardAppendsEventToGameEvents(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	initialEventCount := len(g.Events)

	events := DrawCard(g, 0)

	if len(events) != 1 {
		t.Fatalf("expected 1 returned event, got %d", len(events))
	}

	if len(g.Events) != initialEventCount+1 {
		t.Fatalf("expected game events size %d, got %d", initialEventCount+1, len(g.Events))
	}

	lastEvent := g.Events[len(g.Events)-1]

	if lastEvent.Type != EventTypeCardDrawn {
		t.Fatalf("expected last game event type %q, got %q", EventTypeCardDrawn, lastEvent.Type)
	}
}

// TestDrawCardDoesNotPanicWithEmptyDeck verifies that drawing from an empty deck
// does not mutate the player's hand or deck and still returns a safe event.
func TestDrawCardDoesNotPanicWithEmptyDeck(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	player := &g.Players[0]
	player.Deck = []CardInstance{}

	initialHandSize := len(player.Hand)

	events := DrawCard(g, 0)

	if len(player.Deck) != 0 {
		t.Fatalf("expected deck to stay empty, got %d cards", len(player.Deck))
	}

	if len(player.Hand) != initialHandSize {
		t.Fatalf("expected hand size to stay %d, got %d", initialHandSize, len(player.Hand))
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event for empty deck, got %d", len(events))
	}

	if events[0].Type != EventTypeCardDrawn {
		t.Fatalf("expected event type %q, got %q", EventTypeCardDrawn, events[0].Type)
	}
}

// TestDrawCardDoesNotPanicWithInvalidPlayerIndex verifies that invalid player
// indexes are handled safely without panicking or mutating game state.
func TestDrawCardDoesNotPanicWithInvalidPlayerIndex(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	events := DrawCard(g, -1)

	if len(events) != 1 {
		t.Fatalf("expected 1 event for invalid player index, got %d", len(events))
	}

	events = DrawCard(g, 999)

	if len(events) != 1 {
		t.Fatalf("expected 1 event for invalid player index, got %d", len(events))
	}
}

// TestDrawCardDoesNotPanicWithNilGame verifies that DrawCard handles a nil game
// pointer safely and returns an event instead of panicking.
func TestDrawCardDoesNotPanicWithNilGame(t *testing.T) {
	events := DrawCard(nil, 0)

	if len(events) != 1 {
		t.Fatalf("expected 1 event for nil game, got %d", len(events))
	}

	if events[0].Type != EventTypeCardDrawn {
		t.Fatalf("expected event type %q, got %q", EventTypeCardDrawn, events[0].Type)
	}
}
