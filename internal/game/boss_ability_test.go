package game

import "testing"

func TestApplyActionEndTurnAppliesOneBossAbilityEvent(t *testing.T) {
	g := newTestGame()

	events, err := ApplyAction(g, Action{
		Type:     ActionTypeEndTurn,
		PlayerID: g.Players[0].ID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if countEventsByType(events, EventBossAbility) != 1 {
		t.Fatalf("expected exactly 1 %q event, got %d", EventBossAbility, countEventsByType(events, EventBossAbility))
	}

	if countEventsByType(g.Events, EventBossAbility) != 1 {
		t.Fatalf("expected game events to contain exactly 1 %q event, got %d", EventBossAbility, countEventsByType(g.Events, EventBossAbility))
	}
}

func TestApplyBossAbilityZapHeroesDamagesBothHeroes(t *testing.T) {
	g := newTestGame()
	g.Players[0].Health = 10
	g.Players[1].Health = 9

	events := ApplyBossAbility(g, BossAbility{
		Type: BossAbilityZapHeroes,
		Name: "Zap Heroes",
	})

	if g.Players[0].Health != 8 {
		t.Fatalf("expected player 1 health 8, got %d", g.Players[0].Health)
	}

	if g.Players[1].Health != 7 {
		t.Fatalf("expected player 2 health 7, got %d", g.Players[1].Health)
	}

	if countEventsByType(events, EventBossAbility) != 1 {
		t.Fatalf("expected exactly 1 %q event, got %d", EventBossAbility, countEventsByType(events, EventBossAbility))
	}

	if countEventsByType(events, EventDamage) != 2 {
		t.Fatalf("expected exactly 2 %q events, got %d", EventDamage, countEventsByType(events, EventDamage))
	}
}

func TestApplyBossAbilityBombSalvoDamagesOneValidTarget(t *testing.T) {
	g := newTestGame()
	g.Seed = 7
	g.Turn = 3
	g.Boss.Health = 40
	g.Players[0].Health = 10
	g.Players[1].Health = 10
	g.Players[0].Board = []Minion{
		{ID: MinionID("p1_m1"), OwnerID: g.Players[0].ID, Name: "P1M1", Health: 4, MaxHealth: 4},
	}
	g.Players[1].Board = []Minion{
		{ID: MinionID("p2_m1"), OwnerID: g.Players[1].ID, Name: "P2M1", Health: 5, MaxHealth: 5},
	}

	player1Before := g.Players[0].Health
	player2Before := g.Players[1].Health
	p1m1Before := g.Players[0].Board[0].Health
	p2m1Before := g.Players[1].Board[0].Health

	events := ApplyBossAbility(g, BossAbility{
		Type: BossAbilityBombSalvo,
		Name: "Bomb Salvo",
	})

	changed := 0

	if g.Players[0].Health == player1Before-2 {
		changed++
	}
	if g.Players[1].Health == player2Before-2 {
		changed++
	}
	if g.Players[0].Board[0].Health == p1m1Before-2 {
		changed++
	}
	if g.Players[1].Board[0].Health == p2m1Before-2 {
		changed++
	}

	if changed != 1 {
		t.Fatalf("expected exactly one bomb target to take 2 damage, got %d damaged targets", changed)
	}

	if g.Boss.Health != 40 {
		t.Fatalf("expected boss health to stay unchanged, got %d", g.Boss.Health)
	}

	if countEventsByType(events, EventBossAbility) != 1 {
		t.Fatalf("expected exactly 1 %q event, got %d", EventBossAbility, countEventsByType(events, EventBossAbility))
	}
}

func TestApplyBossAbilityOverclockIncreasesBossAttack(t *testing.T) {
	g := newTestGame()
	initialAttack := g.Boss.Attack

	events := ApplyBossAbility(g, BossAbility{
		Type: BossAbilityOverclock,
		Name: "Overclock",
	})

	if g.Boss.Attack != initialAttack+1 {
		t.Fatalf("expected boss attack %d, got %d", initialAttack+1, g.Boss.Attack)
	}

	if countEventsByType(events, EventBossAbility) != 1 {
		t.Fatalf("expected exactly 1 %q event, got %d", EventBossAbility, countEventsByType(events, EventBossAbility))
	}
}

func TestResolveBossAbilityDeterministicForSameSeed(t *testing.T) {
	g1 := newTestGame()
	g2 := newTestGame()

	g1.Seed = 999
	g2.Seed = 999
	g1.Turn = 4
	g2.Turn = 4

	g1.Players[0].Health = 20
	g2.Players[0].Health = 20
	g1.Players[1].Health = 18
	g2.Players[1].Health = 18
	g1.Players[0].Board = []Minion{
		{ID: MinionID("a1"), OwnerID: g1.Players[0].ID, Name: "A1", Health: 4, MaxHealth: 4},
	}
	g2.Players[0].Board = []Minion{
		{ID: MinionID("a1"), OwnerID: g2.Players[0].ID, Name: "A1", Health: 4, MaxHealth: 4},
	}
	g1.Players[1].Board = []Minion{
		{ID: MinionID("b1"), OwnerID: g1.Players[1].ID, Name: "B1", Health: 3, MaxHealth: 3},
	}
	g2.Players[1].Board = []Minion{
		{ID: MinionID("b1"), OwnerID: g2.Players[1].ID, Name: "B1", Health: 3, MaxHealth: 3},
	}

	events1 := ResolveBossAbility(g1)
	events2 := ResolveBossAbility(g2)

	if len(events1) == 0 || len(events2) == 0 {
		t.Fatal("expected boss ability events to be returned")
	}

	if events1[0].Type != EventBossAbility || events2[0].Type != EventBossAbility {
		t.Fatalf("expected first event type %q in both runs", EventBossAbility)
	}

	if events1[0].Message != events2[0].Message {
		t.Fatalf("expected same boss ability, got %q and %q", events1[0].Message, events2[0].Message)
	}

	if g1.Boss.Attack != g2.Boss.Attack {
		t.Fatalf("expected equal boss attack, got %d and %d", g1.Boss.Attack, g2.Boss.Attack)
	}

	if g1.Players[0].Health != g2.Players[0].Health || g1.Players[1].Health != g2.Players[1].Health {
		t.Fatalf(
			"expected equal hero healths, got p1:%d/%d p2:%d/%d",
			g1.Players[0].Health,
			g2.Players[0].Health,
			g1.Players[1].Health,
			g2.Players[1].Health,
		)
	}

	if g1.Players[0].Board[0].Health != g2.Players[0].Board[0].Health {
		t.Fatalf("expected equal minion health for player 1, got %d and %d", g1.Players[0].Board[0].Health, g2.Players[0].Board[0].Health)
	}

	if g1.Players[1].Board[0].Health != g2.Players[1].Board[0].Health {
		t.Fatalf("expected equal minion health for player 2, got %d and %d", g1.Players[1].Board[0].Health, g2.Players[1].Board[0].Health)
	}
}

func countEventsByType(events []GameEvent, eventType EventType) int {
	count := 0

	for i := range events {
		if events[i].Type == eventType {
			count++
		}
	}

	return count
}
