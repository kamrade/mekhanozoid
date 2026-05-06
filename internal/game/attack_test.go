package game

import (
	"errors"
	"testing"
)

func TestApplyActionAttackDamagesBossAndExhaustsMinion(t *testing.T) {
	g := newTestGame()
	player := &g.Players[0]
	player.Board = []Minion{
		{
			ID:        MinionID("player_1_drone_1"),
			OwnerID:   player.ID,
			Name:      "Drone",
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}

	initialBossHealth := g.Boss.Health

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Boss.Health != initialBossHealth-2 {
		t.Fatalf("expected boss health %d, got %d", initialBossHealth-2, g.Boss.Health)
	}

	if g.Players[0].Board[0].CanAttack {
		t.Fatal("expected minion CanAttack=false after attack")
	}

	if !hasEventType(events, EventAttack) {
		t.Fatalf("expected returned events to contain %q", EventAttack)
	}

	if !hasEventType(events, EventDamage) {
		t.Fatalf("expected returned events to contain %q", EventDamage)
	}
}

func TestApplyActionAttackRejectsMinionCantAttack(t *testing.T) {
	g := newTestGame()
	player := &g.Players[0]
	player.Board = []Minion{
		{
			ID:        MinionID("player_1_drone_1"),
			OwnerID:   player.ID,
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: false,
		},
	}

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if !errors.Is(err, ErrMinionCantAttack) {
		t.Fatalf("expected ErrMinionCantAttack, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionAttackRejectsSecondAttackSameTurn(t *testing.T) {
	g := newTestGame()
	player := &g.Players[0]
	player.Board = []Minion{
		{
			ID:        MinionID("player_1_drone_1"),
			OwnerID:   player.ID,
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if err != nil {
		t.Fatalf("expected first attack to succeed, got %v", err)
	}

	snapshot := snapshotGameState(g)

	_, err = ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if !errors.Is(err, ErrMinionCantAttack) {
		t.Fatalf("expected ErrMinionCantAttack, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionAttackRejectsInactivePlayer(t *testing.T) {
	g := newTestGame()
	inactivePlayer := &g.Players[1]
	inactivePlayer.Board = []Minion{
		{
			ID:        MinionID("player_2_drone_1"),
			OwnerID:   inactivePlayer.ID,
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: inactivePlayer.ID,
		SourceID: inactivePlayer.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if !errors.Is(err, ErrNotYourTurn) {
		t.Fatalf("expected ErrNotYourTurn, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionAttackRejectsAttackingWithOtherPlayersMinion(t *testing.T) {
	g := newTestGame()
	activePlayer := &g.Players[0]
	otherPlayer := &g.Players[1]
	otherPlayer.Board = []Minion{
		{
			ID:        MinionID("player_2_drone_1"),
			OwnerID:   otherPlayer.ID,
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: activePlayer.ID,
		SourceID: otherPlayer.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if !errors.Is(err, ErrMinionNotFound) {
		t.Fatalf("expected ErrMinionNotFound, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionAttackRejectsInvalidTarget(t *testing.T) {
	g := newTestGame()
	player := &g.Players[0]
	player.Board = []Minion{
		{
			ID:        MinionID("player_1_drone_1"),
			OwnerID:   player.ID,
			Attack:    2,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDHero0,
	})
	if !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("expected ErrInvalidTarget, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionAttackCanTriggerGameWon(t *testing.T) {
	g := newTestGame()
	player := &g.Players[0]
	player.Board = []Minion{
		{
			ID:        MinionID("player_1_drone_1"),
			OwnerID:   player.ID,
			Attack:    5,
			Health:    3,
			MaxHealth: 3,
			CanAttack: true,
		},
	}
	g.Boss.Health = 4

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeAttack,
		PlayerID: player.ID,
		SourceID: player.Board[0].ID,
		TargetID: TargetIDBoss,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Status != GameStatusWon {
		t.Fatalf("expected status %q, got %q", GameStatusWon, g.Status)
	}

	if g.Boss.Health != 0 {
		t.Fatalf("expected boss health 0, got %d", g.Boss.Health)
	}

	if !hasEventType(events, EventGameWon) {
		t.Fatalf("expected returned events to contain %q", EventGameWon)
	}
}
