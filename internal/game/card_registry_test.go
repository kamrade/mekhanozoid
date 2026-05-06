package game

import "testing"

func TestCardRegistryIsNotEmpty(t *testing.T) {
	if len(CardRegistry) == 0 {
		t.Fatal("expected card registry not to be empty")
	}
}

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

func TestCardRegistryDefinitionsAreValid(t *testing.T) {
	if err := ValidateCardRegistry(); err != nil {
		t.Fatal(err)
	}
}

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

func TestDroneStats(t *testing.T) {
	drone := CardRegistry[CardID("drone")]

	if drone.Attack != 2 {
		t.Fatalf("expected drone attack 2, got %d", drone.Attack)
	}

	if drone.Health != 3 {
		t.Fatalf("expected drone health 3, got %d", drone.Health)
	}
}
