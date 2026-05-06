package game

import (
	"errors"
	"testing"
)

// TestApplyActionPlayCardRepairHealsHero0 verifies that Repair can heal the first hero.
func TestApplyActionPlayCardRepairHealsHero0(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 2
	player.MaxMana = 2
	g.Players[0].Health = 20

	initialHandSize := len(player.Hand)
	initialEventCount := len(g.Events)

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: TargetIDHero0,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Players[0].Health != 25 {
		t.Fatalf("expected hero 0 health 25, got %d", g.Players[0].Health)
	}

	if player.Mana != 0 {
		t.Fatalf("expected mana 0, got %d", player.Mana)
	}

	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("expected hand size %d, got %d", initialHandSize-1, len(player.Hand))
	}

	if hasCardInHand(player, card.ID) {
		t.Fatalf("expected repair card %q to be removed from hand", card.ID)
	}

	if !hasEventType(events, EventCardPlayed) {
		t.Fatalf("expected returned events to contain %q", EventCardPlayed)
	}

	if !hasEventType(events, EventHeal) {
		t.Fatalf("expected returned events to contain %q", EventHeal)
	}

	if len(g.Events) != initialEventCount+len(events) {
		t.Fatalf("expected game events size %d, got %d", initialEventCount+len(events), len(g.Events))
	}
}

// TestApplyActionPlayCardRepairHealsHero1 verifies that Repair can heal the second hero.
func TestApplyActionPlayCardRepairHealsHero1(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 2
	player.MaxMana = 2
	g.Players[1].Health = 10

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: TargetIDHero1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Players[1].Health != 15 {
		t.Fatalf("expected hero 1 health 15, got %d", g.Players[1].Health)
	}

	if !hasEventType(events, EventHeal) {
		t.Fatalf("expected returned events to contain %q", EventHeal)
	}
}

// TestApplyActionPlayCardRepairRequiresTarget verifies that Repair cannot be played without a target.
func TestApplyActionPlayCardRepairRequiresTarget(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 2
	player.MaxMana = 2

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if !errors.Is(err, ErrTargetRequired) {
		t.Fatalf("expected ErrTargetRequired, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

// TestApplyActionPlayCardRepairRejectsBossTarget verifies that Repair cannot target the boss.
func TestApplyActionPlayCardRepairRejectsBossTarget(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 2
	player.MaxMana = 2

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: TargetIDBoss,
	})

	if !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("expected ErrInvalidTarget, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

// TestApplyActionPlayCardRepairRejectsUnknownTarget verifies that Repair rejects unknown target IDs.
func TestApplyActionPlayCardRepairRejectsUnknownTarget(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 2
	player.MaxMana = 2

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: "hero:99",
	})

	if !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("expected ErrInvalidTarget, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

// TestApplyActionPlayCardRepairDoesNotExceedMaxHealth verifies that healing is capped by MaxHealth.
func TestApplyActionPlayCardRepairDoesNotExceedMaxHealth(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 2
	player.MaxMana = 2
	g.Players[0].Health = g.Players[0].MaxHealth - 2

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: TargetIDHero0,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Players[0].Health != g.Players[0].MaxHealth {
		t.Fatalf("expected hero 0 health to be capped at %d, got %d", g.Players[0].MaxHealth, g.Players[0].Health)
	}

	healEvent := findFirstEvent(events, EventHeal)
	if healEvent == nil {
		t.Fatalf("expected returned events to contain %q", EventHeal)
	}

	if healEvent.Amount != 2 {
		t.Fatalf("expected actual healed amount 2, got %d", healEvent.Amount)
	}
}

// TestApplyActionPlayCardRepairRejectsNotEnoughMana verifies that Repair cannot be played without enough mana.
func TestApplyActionPlayCardRepairRejectsNotEnoughMana(t *testing.T) {
	g := newTestGameWithRepairInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 1
	player.MaxMana = 2

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
		TargetID: TargetIDHero0,
	})

	if !errors.Is(err, ErrNotEnoughMana) {
		t.Fatalf("expected ErrNotEnoughMana, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func newTestGameWithRepairInPlayer1Hand() *Game {
	g := newTestGame()

	player := &g.Players[0]
	player.Hand = []CardInstance{
		{
			ID:           CardInstanceID("player_1_repair_1"),
			DefinitionID: CardID("repair"),
			OwnerID:      player.ID,
		},
	}

	return g
}

func findFirstEvent(events []GameEvent, eventType EventType) *GameEvent {
	for i := range events {
		if events[i].Type == eventType {
			return &events[i]
		}
	}

	return nil
}
