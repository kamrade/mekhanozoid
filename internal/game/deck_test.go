package game

import "testing"

// TestStartingDeckSizeMatchesCardIDs verifies that the declared starting deck size
// matches the actual number of card IDs used to build starting decks.
func TestStartingDeckSizeMatchesCardIDs(t *testing.T) {
	if StartingDeckSize != len(StartingDeckCardIDs) {
		t.Fatalf(
			"expected StartingDeckSize %d to match len(StartingDeckCardIDs) %d",
			StartingDeckSize,
			len(StartingDeckCardIDs),
		)
	}
}

// TestStartingDeckCardsExistInRegistry verifies that every card in a generated
// starting deck references a known card definition from CardRegistry.
func TestStartingDeckCardsExistInRegistry(t *testing.T) {
	deck := NewStartingDeck(PlayerID("player_1"), 42)

	if err := ValidateDeckCardsExist(deck); err != nil {
		t.Fatal(err)
	}
}

// TestNewGameHasNoUnknownCards verifies that NewGame creates players whose
// decks and hands only contain card instances backed by CardRegistry definitions.
func TestNewGameHasNoUnknownCards(t *testing.T) {
	g := NewGame(
		"game_1",
		PlayerConfig{ID: "player_1", Name: "Player 1"},
		PlayerConfig{ID: "player_2", Name: "Player 2"},
		42,
	)

	for _, player := range g.Players {
		if err := ValidateDeckCardsExist(player.Deck); err != nil {
			t.Fatalf("player %s deck has unknown card: %v", player.ID, err)
		}

		if err := ValidateDeckCardsExist(player.Hand); err != nil {
			t.Fatalf("player %s hand has unknown card: %v", player.ID, err)
		}
	}
}
