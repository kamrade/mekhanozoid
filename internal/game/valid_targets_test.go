package game

import (
	"errors"
	"testing"
)

// TestValidTargetsForStrikeReturnsEmptyList verifies that Strike does not require
// manual target selection and therefore returns no selectable targets.
func TestValidTargetsForStrikeReturnsEmptyList(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := g.Players[0]
	card := player.Hand[0]

	targets, err := ValidTargets(g, string(player.ID), string(card.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(targets) != 0 {
		t.Fatalf("expected no valid targets for strike, got %d", len(targets))
	}
}

// TestValidTargetsForRepairReturnsBothHeroes verifies that Repair can target
// both heroes and only heroes.
func TestValidTargetsForRepairReturnsBothHeroes(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := g.Players[0]
	card := player.Hand[0]

	targets, err := ValidTargets(g, string(player.ID), string(card.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("expected 2 valid targets for repair, got %d", len(targets))
	}

	if !hasTargetID(targets, TargetIDHero0) {
		t.Fatalf("expected targets to contain %q", TargetIDHero0)
	}

	if !hasTargetID(targets, TargetIDHero1) {
		t.Fatalf("expected targets to contain %q", TargetIDHero1)
	}

	if hasTargetID(targets, TargetIDBoss) {
		t.Fatalf("expected targets not to contain %q", TargetIDBoss)
	}
}

// TestValidTargetsForRepairIncludesDisplayData verifies that valid targets
// include enough metadata for a future UI to render target choices.
func TestValidTargetsForRepairIncludesDisplayData(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := g.Players[0]
	card := player.Hand[0]

	targets, err := ValidTargets(g, string(player.ID), string(card.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, target := range targets {
		if target.ID == "" {
			t.Fatal("expected target ID to be set")
		}

		if target.Kind != TargetKindHero {
			t.Fatalf("expected target kind %q, got %q", TargetKindHero, target.Kind)
		}

		if target.Type != TargetTypePlayer {
			t.Fatalf("expected target type %q, got %q", TargetTypePlayer, target.Type)
		}

		if target.PlayerID == "" {
			t.Fatal("expected target player ID to be set")
		}

		if target.DisplayName == "" {
			t.Fatal("expected target display name to be set")
		}
	}
}

// TestValidTargetsReturnsErrorForUnknownPlayer verifies that ValidTargets
// rejects player IDs that do not exist in the game.
func TestValidTargetsReturnsErrorForUnknownPlayer(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	card := g.Players[0].Hand[0]

	_, err := ValidTargets(g, "missing_player", string(card.ID))

	if !errors.Is(err, ErrInvalidPlayerIndex) {
		t.Fatalf("expected ErrInvalidPlayerIndex, got %v", err)
	}
}

// TestValidTargetsReturnsErrorForCardNotInHand verifies that ValidTargets
// rejects card instance IDs that are not in the selected player's hand.
func TestValidTargetsReturnsErrorForCardNotInHand(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := g.Players[0]

	_, err := ValidTargets(g, string(player.ID), "missing_card")

	if !errors.Is(err, ErrCardNotInHand) {
		t.Fatalf("expected ErrCardNotInHand, got %v", err)
	}
}

// TestValidTargetsReturnsErrorForUnknownCardDefinition verifies that ValidTargets
// rejects card instances whose definition is missing from CardRegistry.
func TestValidTargetsReturnsErrorForUnknownCardDefinition(t *testing.T) {
	g := newTestGame()

	player := &g.Players[0]
	player.Hand = []CardInstance{
		{
			ID:           CardInstanceID("unknown_card_instance"),
			DefinitionID: CardID("unknown_card_definition"),
			OwnerID:      player.ID,
		},
	}

	_, err := ValidTargets(g, string(player.ID), string(player.Hand[0].ID))

	if !errors.Is(err, ErrUnknownCard) {
		t.Fatalf("expected ErrUnknownCard, got %v", err)
	}
}

// TestValidTargetsReturnsErrorForNilGame verifies that ValidTargets handles a nil
// game pointer safely.
func TestValidTargetsReturnsErrorForNilGame(t *testing.T) {
	_, err := ValidTargets(nil, "player_1", "card_1")

	if !errors.Is(err, ErrNilGame) {
		t.Fatalf("expected ErrNilGame, got %v", err)
	}
}

func hasTargetID(targets []Target, targetID string) bool {
	for _, target := range targets {
		if target.ID == targetID {
			return true
		}
	}

	return false
}
