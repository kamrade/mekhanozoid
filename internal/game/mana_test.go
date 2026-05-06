package game

import (
	"errors"
	"testing"
)

// TestNewGameCreatesPlayersWithInitialMana verifies that players created by
// NewGame start with the configured initial mana and max mana values.
func TestNewGameCreatesPlayersWithInitialMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	for _, player := range g.Players {
		if player.Mana != StartingMana {
			t.Fatalf("expected player %s mana %d, got %d", player.ID, StartingMana, player.Mana)
		}

		if player.MaxMana != StartingMaxMana {
			t.Fatalf("expected player %s max mana %d, got %d", player.ID, StartingMaxMana, player.MaxMana)
		}
	}
}

// TestRefreshManaIncreasesMaxManaAndRefillsMana verifies that RefreshMana
// increases a player's max mana by one and refills current mana to that value.
func TestRefreshManaIncreasesMaxManaAndRefillsMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	RefreshMana(g, 0)

	player := g.Players[0]

	if player.MaxMana != 1 {
		t.Fatalf("expected max mana 1, got %d", player.MaxMana)
	}

	if player.Mana != player.MaxMana {
		t.Fatalf("expected mana to equal max mana, got mana=%d max=%d", player.Mana, player.MaxMana)
	}
}

// TestRefreshManaDoesNotExceedMaxMana verifies that repeated mana refreshes
// never increase a player's max mana beyond the global MaxMana cap.
func TestRefreshManaDoesNotExceedMaxMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	for i := 0; i < 20; i++ {
		RefreshMana(g, 0)
	}

	player := g.Players[0]

	if player.MaxMana != MaxMana {
		t.Fatalf("expected max mana %d, got %d", MaxMana, player.MaxMana)
	}

	if player.Mana != MaxMana {
		t.Fatalf("expected mana %d, got %d", MaxMana, player.Mana)
	}
}

// TestSpendManaSubtractsAvailableMana verifies that SpendMana subtracts the
// requested amount when the player has enough available mana.
func TestSpendManaSubtractsAvailableMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	player := &g.Players[0]
	player.Mana = 5
	player.MaxMana = 5

	err := SpendMana(g, 0, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player.Mana != 2 {
		t.Fatalf("expected mana 2, got %d", player.Mana)
	}
}

// TestSpendManaCanSpendAllMana verifies that SpendMana can reduce a player's
// current mana to zero when spending exactly the available amount.
func TestSpendManaCanSpendAllMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	player := &g.Players[0]
	player.Mana = 3
	player.MaxMana = 3

	err := SpendMana(g, 0, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player.Mana != 0 {
		t.Fatalf("expected mana 0, got %d", player.Mana)
	}
}

// TestSpendManaReturnsErrorWhenNotEnoughMana verifies that SpendMana rejects
// spending more mana than the player has and leaves the player's mana unchanged.
func TestSpendManaReturnsErrorWhenNotEnoughMana(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	player := &g.Players[0]
	player.Mana = 2
	player.MaxMana = 2

	err := SpendMana(g, 0, 3)

	if !errors.Is(err, ErrNotEnoughMana) {
		t.Fatalf("expected ErrNotEnoughMana, got %v", err)
	}

	if player.Mana != 2 {
		t.Fatalf("expected mana to remain 2 after error, got %d", player.Mana)
	}
}

// TestSpendManaReturnsErrorForNegativeAmount verifies that SpendMana rejects
// negative amounts and does not mutate the player's current mana.
func TestSpendManaReturnsErrorForNegativeAmount(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	player := &g.Players[0]
	player.Mana = 5
	player.MaxMana = 5

	err := SpendMana(g, 0, -1)

	if !errors.Is(err, ErrNegativeManaAmount) {
		t.Fatalf("expected ErrNegativeManaAmount, got %v", err)
	}

	if player.Mana != 5 {
		t.Fatalf("expected mana to remain 5 after error, got %d", player.Mana)
	}
}

// TestSpendManaReturnsErrorForInvalidPlayerIndex verifies that SpendMana rejects
// indexes that do not point to an existing player.
func TestSpendManaReturnsErrorForInvalidPlayerIndex(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	err := SpendMana(g, -1, 1)

	if !errors.Is(err, ErrInvalidPlayerIndex) {
		t.Fatalf("expected ErrInvalidPlayerIndex, got %v", err)
	}

	err = SpendMana(g, 99, 1)

	if !errors.Is(err, ErrInvalidPlayerIndex) {
		t.Fatalf("expected ErrInvalidPlayerIndex, got %v", err)
	}
}

// TestSpendManaReturnsErrorForNilGame verifies that SpendMana handles a nil game
// pointer safely and returns ErrNilGame.
func TestSpendManaReturnsErrorForNilGame(t *testing.T) {
	err := SpendMana(nil, 0, 1)

	if !errors.Is(err, ErrNilGame) {
		t.Fatalf("expected ErrNilGame, got %v", err)
	}
}

// TestRefreshManaDoesNotPanicForInvalidState verifies that RefreshMana safely
// ignores nil games and invalid player indexes.
func TestRefreshManaDoesNotPanicForInvalidState(t *testing.T) {
	RefreshMana(nil, 0)

	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	RefreshMana(g, -1)
	RefreshMana(g, 99)
}
