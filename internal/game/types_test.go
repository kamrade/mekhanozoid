package game

import "testing"

// TestNewPlayer verifies that NewPlayer initializes the player's identity,
// starting health, and empty card zones.
func TestNewPlayer(t *testing.T) {
	player := NewPlayer(PlayerID("player_1"), "Player 1")

	if player.ID != PlayerID("player_1") {
		t.Fatalf("expected player ID %q, got %q", PlayerID("player_1"), player.ID)
	}

	if player.Name != "Player 1" {
		t.Fatalf("expected player name %q, got %q", "Player 1", player.Name)
	}

	if player.Health <= 0 {
		t.Fatalf("expected positive player health, got %d", player.Health)
	}

	if player.Deck == nil {
		t.Fatal("expected deck to be initialized")
	}

	if player.Hand == nil {
		t.Fatal("expected hand to be initialized")
	}

	if player.Board == nil {
		t.Fatal("expected board to be initialized")
	}

	if player.Discard == nil {
		t.Fatal("expected discard to be initialized")
	}
}

// TestNewBoss verifies that NewBoss initializes the boss identity, name,
// health, and max health values.
func TestNewBoss(t *testing.T) {
	boss := NewBoss(BossID("boss_1"), "Mechazod")

	if boss.ID != BossID("boss_1") {
		t.Fatalf("expected boss ID %q, got %q", BossID("boss_1"), boss.ID)
	}

	if boss.Name != "Mechazod" {
		t.Fatalf("expected boss name %q, got %q", "Mechazod", boss.Name)
	}

	if boss.Health <= 0 {
		t.Fatalf("expected positive boss health, got %d", boss.Health)
	}

	if boss.MaxHealth < boss.Health {
		t.Fatalf("expected max health >= health, got max=%d health=%d", boss.MaxHealth, boss.Health)
	}
}

// TestCardDefinition verifies that a card definition can describe a valid card
// template with an ID and card type.
func TestCardDefinition(t *testing.T) {
	card := CardDefinition{
		ID:          CardID("strike"),
		Name:        "Strike",
		Type:        CardTypeSpell,
		Cost:        1,
		Description: "Deal damage.",
	}

	if card.ID == "" {
		t.Fatal("expected card definition ID to be set")
	}

	if card.Type != CardTypeSpell {
		t.Fatalf("expected card type %q, got %q", CardTypeSpell, card.Type)
	}
}

// TestCardInstance verifies that a card instance references a card definition
// and belongs to a specific owner.
func TestCardInstance(t *testing.T) {
	card := CardInstance{
		ID:           CardInstanceID("card_instance_1"),
		DefinitionID: CardID("strike"),
		OwnerID:      PlayerID("player_1"),
	}

	if card.ID == "" {
		t.Fatal("expected card instance ID to be set")
	}

	if card.DefinitionID == "" {
		t.Fatal("expected card definition ID to be set")
	}

	if card.OwnerID == "" {
		t.Fatal("expected owner ID to be set")
	}
}

// TestMinion verifies that a minion can represent a summoned board entity with
// identity, ownership, combat stats, and health values.
func TestMinion(t *testing.T) {
	minion := Minion{
		ID:           MinionID("minion_1"),
		DefinitionID: CardID("bot"),
		OwnerID:      PlayerID("player_1"),
		Name:         "Bot",
		Attack:       2,
		Health:       3,
		MaxHealth:    3,
		Exhausted:    true,
	}

	if minion.ID == "" {
		t.Fatal("expected minion ID to be set")
	}

	if minion.Health <= 0 {
		t.Fatalf("expected positive minion health, got %d", minion.Health)
	}

	if minion.MaxHealth < minion.Health {
		t.Fatalf("expected max health >= health, got max=%d health=%d", minion.MaxHealth, minion.Health)
	}
}

// TestAction verifies that an action can describe player intent, including the
// action type, acting player, card involved, and target.
func TestAction(t *testing.T) {
	action := Action{
		Type:     ActionTypePlayCard,
		PlayerID: PlayerID("player_1"),
		CardID:   CardInstanceID("card_instance_1"),
		Target: Target{
			Type:   TargetTypeBoss,
			BossID: BossID("boss_1"),
		},
	}

	if action.Type != ActionTypePlayCard {
		t.Fatalf("expected action type %q, got %q", ActionTypePlayCard, action.Type)
	}

	if action.Target.Type != TargetTypeBoss {
		t.Fatalf("expected target type %q, got %q", TargetTypeBoss, action.Target.Type)
	}
}

// TestGameEvent verifies that a game event can describe a state change with
// event type, involved player, target, amount, message, and turn number.
func TestGameEvent(t *testing.T) {
	event := GameEvent{
		Type:     EventTypeDamageDealt,
		PlayerID: PlayerID("player_1"),
		Target: Target{
			Type:   TargetTypeBoss,
			BossID: BossID("boss_1"),
		},
		Amount:  3,
		Message: "Player dealt damage to boss.",
		Turn:    1,
	}

	if event.Type != EventTypeDamageDealt {
		t.Fatalf("expected event type %q, got %q", EventTypeDamageDealt, event.Type)
	}

	if event.Amount != 3 {
		t.Fatalf("expected amount 3, got %d", event.Amount)
	}

	if event.Turn != 1 {
		t.Fatalf("expected turn 1, got %d", event.Turn)
	}
}
