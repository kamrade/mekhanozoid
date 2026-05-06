package game

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame()

	if g == nil {
		t.Fatal("expected game to be created")
	}

	if g.Status() != StatusCreated {
		t.Fatalf("expected status %q, got %q", StatusCreated, g.Status())
	}
}
