package game

import (
	"errors"
	"testing"
)

func TestApplyActionPlayCardStrikeDamagesBoss(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 1
	player.MaxMana = 1

	initialBossHealth := g.Boss.Health
	initialHandSize := len(player.Hand)
	initialMana := player.Mana
	initialEventCount := len(g.Events)

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player.Mana != initialMana-1 {
		t.Fatalf("expected mana %d, got %d", initialMana-1, player.Mana)
	}

	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("expected hand size %d, got %d", initialHandSize-1, len(player.Hand))
	}

	if hasCardInHand(player, card.ID) {
		t.Fatalf("expected card %q to be removed from hand", card.ID)
	}

	if g.Boss.Health != initialBossHealth-3 {
		t.Fatalf("expected boss health %d, got %d", initialBossHealth-3, g.Boss.Health)
	}

	if !hasEventType(events, EventCardPlayed) {
		t.Fatalf("expected returned events to contain %q", EventCardPlayed)
	}

	if !hasEventType(events, EventDamage) {
		t.Fatalf("expected returned events to contain %q", EventDamage)
	}

	if len(g.Events) != initialEventCount+len(events) {
		t.Fatalf("expected game events size %d, got %d", initialEventCount+len(events), len(g.Events))
	}
}

func TestApplyActionPlayCardRejectsInactivePlayer(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	inactivePlayer := &g.Players[1]
	inactivePlayer.Hand = []CardInstance{
		{
			ID:           CardInstanceID("player_2_strike_1"),
			DefinitionID: CardID("strike"),
			OwnerID:      inactivePlayer.ID,
		},
	}
	inactivePlayer.Mana = 1
	inactivePlayer.MaxMana = 1

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: inactivePlayer.ID,
		CardID:   inactivePlayer.Hand[0].ID,
	})

	if !errors.Is(err, ErrNotYourTurn) {
		t.Fatalf("expected ErrNotYourTurn, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardRejectsCardNotInHand(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	player.Mana = 1
	player.MaxMana = 1

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   CardInstanceID("missing_card"),
	})

	if !errors.Is(err, ErrCardNotInHand) {
		t.Fatalf("expected ErrCardNotInHand, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardRejectsNotEnoughMana(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 0
	player.MaxMana = 1

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if !errors.Is(err, ErrNotEnoughMana) {
		t.Fatalf("expected ErrNotEnoughMana, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardRejectsUnknownCardDefinition(t *testing.T) {
	g := newTestGame()

	player := &g.Players[0]
	player.Hand = []CardInstance{
		{
			ID:           CardInstanceID("unknown_card_instance"),
			DefinitionID: CardID("unknown_card_definition"),
			OwnerID:      player.ID,
		},
	}
	player.Mana = 10
	player.MaxMana = 10

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   player.Hand[0].ID,
	})

	if !errors.Is(err, ErrUnknownCard) {
		t.Fatalf("expected ErrUnknownCard, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardRejectsWhenGameIsNotActive(t *testing.T) {
	g := newTestGameWithStrikeInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]

	player.Mana = 1
	player.MaxMana = 1
	g.Status = GameStatusWon

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})

	if !errors.Is(err, ErrGameNotActive) {
		t.Fatalf("expected ErrGameNotActive, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardDroneSummonsMinion(t *testing.T) {
	g := newTestGameWithDroneInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 2
	player.MaxMana = 2

	initialHandSize := len(player.Hand)
	initialBoardSize := len(player.Board)
	initialEventCount := len(g.Events)

	events, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player.Mana != 0 {
		t.Fatalf("expected mana 0, got %d", player.Mana)
	}

	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("expected hand size %d, got %d", initialHandSize-1, len(player.Hand))
	}

	if len(player.Board) != initialBoardSize+1 {
		t.Fatalf("expected board size %d, got %d", initialBoardSize+1, len(player.Board))
	}

	minion := player.Board[len(player.Board)-1]
	if minion.DefinitionID != CardID("drone") {
		t.Fatalf("expected summoned minion definition %q, got %q", CardID("drone"), minion.DefinitionID)
	}
	if minion.Attack != 2 {
		t.Fatalf("expected summoned minion attack 2, got %d", minion.Attack)
	}
	if minion.Health != 3 {
		t.Fatalf("expected summoned minion health 3, got %d", minion.Health)
	}
	if minion.MaxHealth != 3 {
		t.Fatalf("expected summoned minion max health 3, got %d", minion.MaxHealth)
	}
	if minion.CanAttack {
		t.Fatal("expected summoned minion to not be able to attack this turn")
	}

	if !hasEventType(events, EventCardPlayed) {
		t.Fatalf("expected returned events to contain %q", EventCardPlayed)
	}
	if !hasEventType(events, EventTypeMinionSummoned) {
		t.Fatalf("expected returned events to contain %q", EventTypeMinionSummoned)
	}

	if len(g.Events) != initialEventCount+len(events) {
		t.Fatalf("expected game events size %d, got %d", initialEventCount+len(events), len(g.Events))
	}
}

func TestApplyActionPlayCardDroneRejectsWhenBoardIsFull(t *testing.T) {
	g := newTestGameWithDroneInPlayer1Hand()

	player := &g.Players[0]
	card := player.Hand[0]
	player.Mana = 10
	player.MaxMana = 10
	player.Board = make([]Minion, MaxBoardSize)

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: player.ID,
		CardID:   card.ID,
	})
	if !errors.Is(err, ErrBoardFull) {
		t.Fatalf("expected ErrBoardFull, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func TestApplyActionPlayCardDroneRejectsInactivePlayer(t *testing.T) {
	g := newTestGameWithDroneInPlayer1Hand()

	inactivePlayer := &g.Players[1]
	inactivePlayer.Hand = []CardInstance{
		{
			ID:           CardInstanceID("player_2_drone_1"),
			DefinitionID: CardID("drone"),
			OwnerID:      inactivePlayer.ID,
		},
	}
	inactivePlayer.Mana = 2
	inactivePlayer.MaxMana = 2

	snapshot := snapshotGameState(g)

	_, err := ApplyAction(g, Action{
		Type:     ActionTypePlayCard,
		PlayerID: inactivePlayer.ID,
		CardID:   inactivePlayer.Hand[0].ID,
	})
	if !errors.Is(err, ErrNotYourTurn) {
		t.Fatalf("expected ErrNotYourTurn, got %v", err)
	}

	assertGameStateUnchanged(t, snapshot, g)
}

func newTestGameWithStrikeInPlayer1Hand() *Game {
	g := newTestGame()

	player := &g.Players[0]
	player.Hand = []CardInstance{
		{
			ID:           CardInstanceID("player_1_strike_1"),
			DefinitionID: CardID("strike"),
			OwnerID:      player.ID,
		},
	}

	return g
}

func newTestGameWithDroneInPlayer1Hand() *Game {
	g := newTestGame()

	player := &g.Players[0]
	player.Hand = []CardInstance{
		{
			ID:           CardInstanceID("player_1_drone_1"),
			DefinitionID: CardID("drone"),
			OwnerID:      player.ID,
		},
	}

	return g
}

func hasCardInHand(player *Player, cardID CardInstanceID) bool {
	if player == nil {
		return false
	}

	for _, card := range player.Hand {
		if card.ID == cardID {
			return true
		}
	}

	return false
}

type gameStateSnapshot struct {
	Player1Mana      int
	Player1Health    int
	Player1HandSize  int
	Player1BoardSize int
	Player2Mana      int
	Player2Health    int
	Player2HandSize  int
	Player2BoardSize int
	BossHealth       int
	EventCount       int
	ActivePlayerID   PlayerID
	Turn             int
}

func snapshotGameState(g *Game) gameStateSnapshot {
	return gameStateSnapshot{
		Player1Mana:      g.Players[0].Mana,
		Player1Health:    g.Players[0].Health,
		Player1HandSize:  len(g.Players[0].Hand),
		Player1BoardSize: len(g.Players[0].Board),
		Player2Mana:      g.Players[1].Mana,
		Player2Health:    g.Players[1].Health,
		Player2HandSize:  len(g.Players[1].Hand),
		Player2BoardSize: len(g.Players[1].Board),
		BossHealth:       g.Boss.Health,
		EventCount:       len(g.Events),
		ActivePlayerID:   g.ActivePlayerID,
		Turn:             g.Turn,
	}
}

func assertGameStateUnchanged(t *testing.T, snapshot gameStateSnapshot, g *Game) {
	t.Helper()

	if g.Players[0].Health != snapshot.Player1Health {
		t.Fatalf("expected player 1 health to remain %d, got %d", snapshot.Player1Health, g.Players[0].Health)
	}

	if g.Players[1].Health != snapshot.Player2Health {
		t.Fatalf("expected player 2 health to remain %d, got %d", snapshot.Player2Health, g.Players[1].Health)
	}

	if g.Players[0].Mana != snapshot.Player1Mana {
		t.Fatalf("expected player 1 mana to remain %d, got %d", snapshot.Player1Mana, g.Players[0].Mana)
	}

	if len(g.Players[0].Hand) != snapshot.Player1HandSize {
		t.Fatalf("expected player 1 hand size to remain %d, got %d", snapshot.Player1HandSize, len(g.Players[0].Hand))
	}

	if len(g.Players[0].Board) != snapshot.Player1BoardSize {
		t.Fatalf("expected player 1 board size to remain %d, got %d", snapshot.Player1BoardSize, len(g.Players[0].Board))
	}

	if g.Players[1].Mana != snapshot.Player2Mana {
		t.Fatalf("expected player 2 mana to remain %d, got %d", snapshot.Player2Mana, g.Players[1].Mana)
	}

	if len(g.Players[1].Hand) != snapshot.Player2HandSize {
		t.Fatalf("expected player 2 hand size to remain %d, got %d", snapshot.Player2HandSize, len(g.Players[1].Hand))
	}

	if len(g.Players[1].Board) != snapshot.Player2BoardSize {
		t.Fatalf("expected player 2 board size to remain %d, got %d", snapshot.Player2BoardSize, len(g.Players[1].Board))
	}

	if g.Boss.Health != snapshot.BossHealth {
		t.Fatalf("expected boss health to remain %d, got %d", snapshot.BossHealth, g.Boss.Health)
	}

	if len(g.Events) != snapshot.EventCount {
		t.Fatalf("expected event count to remain %d, got %d", snapshot.EventCount, len(g.Events))
	}

	if g.ActivePlayerID != snapshot.ActivePlayerID {
		t.Fatalf("expected active player to remain %q, got %q", snapshot.ActivePlayerID, g.ActivePlayerID)
	}

	if g.Turn != snapshot.Turn {
		t.Fatalf("expected turn to remain %d, got %d", snapshot.Turn, g.Turn)
	}
}
