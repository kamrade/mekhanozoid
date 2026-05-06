package game

import "testing"

func TestStartingDeckSizeMatchesCardIDs(t *testing.T) {
	if StartingDeckSize != len(StartingDeckCardIDs) {
		t.Fatalf(
			"expected StartingDeckSize %d to match len(StartingDeckCardIDs) %d",
			StartingDeckSize,
			len(StartingDeckCardIDs),
		)
	}
}

func TestStartingDeckCardsExistInRegistry(t *testing.T) {
	deck := NewStartingDeck(PlayerID("player_1"))

	if err := ValidateDeckCardsExist(deck); err != nil {
		t.Fatal(err)
	}
}

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
