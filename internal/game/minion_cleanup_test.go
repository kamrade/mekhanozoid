package game

import "testing"

func TestCleanupDeadMinionsRemovesDeadFromBothBoards(t *testing.T) {
	g := newTestGame()
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_alive"), OwnerID: g.Players[0].ID, Name: "Alive 1", Health: 2, MaxHealth: 2},
		{ID: MinionID("p1_dead"), OwnerID: g.Players[0].ID, Name: "Dead 1", Health: 0, MaxHealth: 2},
	}
	g.Players[1].Board = []Minion{
		{ID: MinionID("p2_dead"), OwnerID: g.Players[1].ID, Name: "Dead 2", Health: -1, MaxHealth: 2},
		{ID: MinionID("p2_alive"), OwnerID: g.Players[1].ID, Name: "Alive 2", Health: 3, MaxHealth: 3},
	}

	events := CleanupDeadMinions(g)

	if len(g.Players[0].Board) != 1 || g.Players[0].Board[0].ID != MinionID("p1_alive") {
		t.Fatalf("expected only p1_alive to remain on player 1 board, got %+v", g.Players[0].Board)
	}

	if len(g.Players[1].Board) != 1 || g.Players[1].Board[0].ID != MinionID("p2_alive") {
		t.Fatalf("expected only p2_alive to remain on player 2 board, got %+v", g.Players[1].Board)
	}

	if countEventsByType(events, EventMinionDied) != 2 {
		t.Fatalf("expected 2 %q events, got %d", EventMinionDied, countEventsByType(events, EventMinionDied))
	}
}

func TestCleanupDeadMinionsDoesNotDuplicateEventsOnRepeatedCalls(t *testing.T) {
	g := newTestGame()
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_dead"), OwnerID: g.Players[0].ID, Name: "Dead 1", Health: 0, MaxHealth: 2},
	}

	events1 := CleanupDeadMinions(g)
	events2 := CleanupDeadMinions(g)

	if countEventsByType(events1, EventMinionDied) != 1 {
		t.Fatalf("expected 1 %q event on first cleanup, got %d", EventMinionDied, countEventsByType(events1, EventMinionDied))
	}

	if countEventsByType(events2, EventMinionDied) != 0 {
		t.Fatalf("expected 0 %q events on second cleanup, got %d", EventMinionDied, countEventsByType(events2, EventMinionDied))
	}
}

func TestApplyBossAbilityBombSalvoTriggersMinionCleanup(t *testing.T) {
	g := newTestGame()
	g.Players[0].Health = 0
	g.Players[1].Health = 0
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_m1"), OwnerID: g.Players[0].ID, Name: "Fragile", Health: 2, MaxHealth: 2},
	}
	g.Players[1].Board = nil

	events := ApplyBossAbility(g, BossAbility{
		Type: BossAbilityBombSalvo,
		Name: "Bomb Salvo",
	})

	if len(g.Players[0].Board) != 0 {
		t.Fatalf("expected dead minion to be cleaned up, got board size %d", len(g.Players[0].Board))
	}

	if countEventsByType(events, EventMinionDied) != 1 {
		t.Fatalf("expected 1 %q event, got %d", EventMinionDied, countEventsByType(events, EventMinionDied))
	}
}

func TestValidTargetsExcludesDeadMinions(t *testing.T) {
	g := newTestGame()
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_alive"), OwnerID: g.Players[0].ID, Name: "Alive", Health: 1, MaxHealth: 2},
		{ID: MinionID("p1_dead"), OwnerID: g.Players[0].ID, Name: "Dead", Health: 0, MaxHealth: 2},
	}

	targets := targetsForRule(g, TargetingRule{
		Required:     true,
		AllowedKinds: []TargetKind{TargetKindMinion},
	})

	if !hasTargetID(targets, minionTargetID(0, MinionID("p1_alive"))) {
		t.Fatal("expected living minion target to be present")
	}

	if hasTargetID(targets, minionTargetID(0, MinionID("p1_dead"))) {
		t.Fatal("expected dead minion target to be excluded")
	}
}

func TestValidateTargetForCardRejectsDeadMinion(t *testing.T) {
	g := newTestGame()
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_dead"), OwnerID: g.Players[0].ID, Name: "Dead", Health: 0, MaxHealth: 2},
	}

	card := CardDefinition{
		ID:   CardID("test_spell"),
		Type: CardTypeSpell,
		Targeting: TargetingRule{
			Required:     true,
			AllowedKinds: []TargetKind{TargetKindMinion},
		},
	}

	_, _, err := ValidateTargetForCard(g, card, minionTargetID(0, MinionID("p1_dead")))
	if err != ErrInvalidTarget {
		t.Fatalf("expected ErrInvalidTarget, got %v", err)
	}
}
