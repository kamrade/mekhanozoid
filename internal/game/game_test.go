package game

import "testing"

// TestNewGameCreatesInitialState verifies that NewGame creates the core game state:
// game ID, active status, two players, active player, boss, and initial turn number.
func TestNewGameCreatesInitialState(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{
			ID:   "player_1",
			Name: "Player 1",
		},
		PlayerConfig{
			ID:   "player_2",
			Name: "Player 2",
		},
		42,
	)

	if g == nil {
		t.Fatal("expected game to be created")
	}

	if g.ID != GameID("game_1") {
		t.Fatalf("expected game ID %q, got %q", GameID("game_1"), g.ID)
	}

	if g.Status != GameStatusActive {
		t.Fatalf("expected status %q, got %q", GameStatusActive, g.Status)
	}

	if len(g.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(g.Players))
	}

	if g.ActivePlayerID == "" {
		t.Fatal("expected active player ID to be set")
	}

	if g.ActivePlayerID != g.Players[0].ID {
		t.Fatalf("expected first player to be active, got %q", g.ActivePlayerID)
	}

	if !g.Players[0].IsCurrent {
		t.Fatal("expected first player to be current")
	}

	if g.Players[1].IsCurrent {
		t.Fatal("expected second player not to be current")
	}

	if g.Boss.ID == "" {
		t.Fatal("expected boss ID to be set")
	}

	if g.Boss.Health <= 0 {
		t.Fatalf("expected boss health to be positive, got %d", g.Boss.Health)
	}

	if g.Turn != 1 {
		t.Fatalf("expected initial turn to be 1, got %d", g.Turn)
	}
}

// TestNewGameCreatesPlayersWithStartingHealth verifies that each player created
// by NewGame starts with the configured starting health.
func TestNewGameCreatesPlayersWithStartingHealth(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	for _, player := range g.Players {
		if player.Health != StartingPlayerHealth {
			t.Fatalf("expected player %s health %d, got %d", player.ID, StartingPlayerHealth, player.Health)
		}
	}
}

// TestNewGameCreatesDecksAndHands verifies that NewGame creates starting decks
// and deals the expected number of cards into each player's hand.
func TestNewGameCreatesDecksAndHands(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	for _, player := range g.Players {
		if len(player.Hand) != StartingHandSize {
			t.Fatalf("expected player %s hand size %d, got %d", player.ID, StartingHandSize, len(player.Hand))
		}

		expectedRemainingDeckSize := StartingDeckSize - StartingHandSize
		if len(player.Deck) != expectedRemainingDeckSize {
			t.Fatalf(
				"expected player %s deck size %d, got %d",
				player.ID,
				expectedRemainingDeckSize,
				len(player.Deck),
			)
		}
	}
}

// TestNewGameIsDeterministicForSameSeed verifies that using the same seed
// creates the same shuffled deck and starting hand order.
func TestNewGameIsDeterministicForSameSeed(t *testing.T) {
	g1 := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	g2 := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	assertSameCards(t, g1.Players[0].Hand, g2.Players[0].Hand)
	assertSameCards(t, g1.Players[0].Deck, g2.Players[0].Deck)
	assertSameCards(t, g1.Players[1].Hand, g2.Players[1].Hand)
	assertSameCards(t, g1.Players[1].Deck, g2.Players[1].Deck)
}

// TestNewGameProducesDifferentShuffleForDifferentSeeds verifies that changing
// the seed changes the resulting card order.
func TestNewGameProducesDifferentShuffleForDifferentSeeds(t *testing.T) {
	g1 := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	g2 := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		43,
	)

	if sameCards(g1.Players[0].Hand, g2.Players[0].Hand) &&
		sameCards(g1.Players[0].Deck, g2.Players[0].Deck) {
		t.Fatal("expected different seeds to produce different player 1 card order")
	}
}

// assertSameCards fails the test if two card slices differ in length or order.
func assertSameCards(t *testing.T, a []CardInstance, b []CardInstance) {
	t.Helper()

	if !sameCards(a, b) {
		t.Fatalf("expected card slices to be equal:\na=%v\nb=%v", a, b)
	}
}

// sameCards reports whether two card slices contain the same cards in the same order.
func sameCards(a []CardInstance, b []CardInstance) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
