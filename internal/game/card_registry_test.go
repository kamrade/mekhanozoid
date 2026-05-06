package game

import "testing"

// TestCardRegistryIsNotEmpty verifies that the card registry contains at least one card.
// An empty registry would make deck creation and card lookup impossible.
func TestCardRegistryIsNotEmpty(t *testing.T) {
	if len(CardRegistry) == 0 {
		t.Fatal("expected card registry not to be empty")
	}
}

// TestCardRegistryHasRequiredCards verifies that the minimal required card set exists.
// These cards are expected by the starting deck and early game prototype.
func TestCardRegistryHasRequiredCards(t *testing.T) {
	requiredCards := []CardID{
		CardID("strike"),
		CardID("repair"),
		CardID("drone"),
	}

	for _, cardID := range requiredCards {
		if _, ok := CardRegistry[cardID]; !ok {
			t.Fatalf("expected card registry to contain %q", cardID)
		}
	}
}

// TestCardRegistryDefinitionsAreValid verifies that all card definitions pass registry validation.
// This catches missing IDs, mismatched registry keys, empty names, and empty card types.
func TestCardRegistryDefinitionsAreValid(t *testing.T) {
	if err := ValidateCardRegistry(); err != nil {
		t.Fatal(err)
	}
}

// TestRequiredCardTypes verifies that the required cards have the expected card types.
// Strike and Repair are spells, while Drone is a minion.
func TestRequiredCardTypes(t *testing.T) {
	strike := CardRegistry[CardID("strike")]
	repair := CardRegistry[CardID("repair")]
	drone := CardRegistry[CardID("drone")]

	if strike.Type != CardTypeSpell {
		t.Fatalf("expected strike to be %q, got %q", CardTypeSpell, strike.Type)
	}

	if repair.Type != CardTypeSpell {
		t.Fatalf("expected repair to be %q, got %q", CardTypeSpell, repair.Type)
	}

	if drone.Type != CardTypeMinion {
		t.Fatalf("expected drone to be %q, got %q", CardTypeMinion, drone.Type)
	}
}

// TestRequiredCardEffects verifies that the required cards describe their intended effects.
// Effects are metadata only at this stage and are not executed by the game engine yet.
func TestRequiredCardEffects(t *testing.T) {
	strike := CardRegistry[CardID("strike")]
	repair := CardRegistry[CardID("repair")]
	drone := CardRegistry[CardID("drone")]

	if strike.Effect.Type != EffectDamageBoss {
		t.Fatalf("expected strike effect %q, got %q", EffectDamageBoss, strike.Effect.Type)
	}

	if strike.Effect.Amount != 3 {
		t.Fatalf("expected strike amount 3, got %d", strike.Effect.Amount)
	}

	if repair.Effect.Type != EffectHealHero {
		t.Fatalf("expected repair effect %q, got %q", EffectHealHero, repair.Effect.Type)
	}

	if repair.Effect.Amount != 3 {
		t.Fatalf("expected repair amount 3, got %d", repair.Effect.Amount)
	}

	if drone.Effect.Type != EffectSummon {
		t.Fatalf("expected drone effect %q, got %q", EffectSummon, drone.Effect.Type)
	}
}

// TestDroneStats verifies that Drone has the expected minion combat stats.
// These stats will later be used when the card is played and a minion is summoned.
func TestDroneStats(t *testing.T) {
	drone := CardRegistry[CardID("drone")]

	if drone.Attack != 2 {
		t.Fatalf("expected drone attack 2, got %d", drone.Attack)
	}

	if drone.Health != 3 {
		t.Fatalf("expected drone health 3, got %d", drone.Health)
	}
}
