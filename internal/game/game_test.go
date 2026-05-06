package game

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame()

	if g == nil {
		t.Fatal("expected game to be created")
	}

	if g.ID == "" {
		t.Fatal("expected game ID to be set")
	}

	if g.Status != GameStatusCreated {
		t.Fatalf("expected status %q, got %q", GameStatusCreated, g.Status)
	}

	if len(g.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(g.Players))
	}

	if g.Boss.ID == "" {
		t.Fatal("expected boss ID to be set")
	}

	if g.Boss.Health <= 0 {
		t.Fatalf("expected boss health to be positive, got %d", g.Boss.Health)
	}

	if g.Turn != 0 {
		t.Fatalf("expected initial turn to be 0, got %d", g.Turn)
	}
}
