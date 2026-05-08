package game

import (
	"fmt"
	"math/rand"
)

type BossAbilityType string

const (
	BossAbilityZapHeroes BossAbilityType = "zap_heroes"
	BossAbilityBombSalvo BossAbilityType = "bomb_salvo"
	BossAbilityOverclock BossAbilityType = "overclock"
)

type BossAbility struct {
	Type BossAbilityType
	Name string
}

func ResolveBossAbility(g *Game) []GameEvent {
	if g == nil {
		return nil
	}

	return ApplyBossAbility(g, ChooseBossAbility(g))
}

func ChooseBossAbility(g *Game) BossAbility {
	abilities := []BossAbility{
		{Type: BossAbilityZapHeroes, Name: "Zap Heroes"},
		{Type: BossAbilityBombSalvo, Name: "Bomb Salvo"},
		{Type: BossAbilityOverclock, Name: "Overclock"},
	}

	if g == nil {
		return abilities[0]
	}

	index := randomIntBySeed(g, 11, len(abilities))
	return abilities[index]
}

func ApplyBossAbility(g *Game, ability BossAbility) []GameEvent {
	if g == nil {
		return nil
	}

	abilityEvent := GameEvent{
		Type:    EventBossAbility,
		Message: fmt.Sprintf("boss used %s", ability.Name),
		Turn:    g.Turn,
		Target: Target{
			Type:        TargetTypeBoss,
			Kind:        TargetKindBoss,
			BossID:      g.Boss.ID,
			DisplayName: g.Boss.Name,
		},
	}

	events := []GameEvent{abilityEvent}

	switch ability.Type {
	case BossAbilityZapHeroes:
		for i := range g.Players {
			events = append(events, damageHero(g, i, 2))
		}
	case BossAbilityBombSalvo:
		targets := collectBombSalvoTargets(g)
		if len(targets) == 0 {
			return events
		}

		targetIndex := randomIntBySeed(g, 29, len(targets))
		selected := targets[targetIndex]

		if selected.playerIndex >= 0 {
			events = append(events, damageHero(g, selected.playerIndex, 2))
		} else {
			events = append(events, damageMinion(g, selected.ownerIndex, selected.minionIndex, 2))
		}
	case BossAbilityOverclock:
		g.Boss.Attack++
	default:
		abilityEvent.Message = "boss used unknown ability"
		events[0] = abilityEvent
	}

	deathEvents := CleanupDeadMinions(g)
	events = append(events, deathEvents...)

	return events
}

type bombSalvoTarget struct {
	playerIndex int
	ownerIndex  int
	minionIndex int
}

func collectBombSalvoTargets(g *Game) []bombSalvoTarget {
	if g == nil {
		return nil
	}

	targets := make([]bombSalvoTarget, 0, len(g.Players))

	for playerIndex := range g.Players {
		if g.Players[playerIndex].Health > 0 {
			targets = append(targets, bombSalvoTarget{playerIndex: playerIndex})
		}

		for minionIndex := range g.Players[playerIndex].Board {
			if g.Players[playerIndex].Board[minionIndex].Health > 0 {
				targets = append(targets, bombSalvoTarget{
					playerIndex: -1,
					ownerIndex:  playerIndex,
					minionIndex: minionIndex,
				})
			}
		}
	}

	return targets
}

func randomIntBySeed(g *Game, salt int64, max int) int {
	if g == nil || max <= 0 {
		return 0
	}

	turnSeed := int64(g.Turn) * 1000003
	r := rand.New(rand.NewSource(g.Seed + turnSeed + salt))
	return r.Intn(max)
}

func damageHero(g *Game, playerIndex int, amount int) GameEvent {
	if g == nil {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage to hero: game is nil",
		}
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage to hero: invalid player index",
			Turn:    g.Turn,
		}
	}

	if amount < 0 {
		amount = 0
	}

	player := &g.Players[playerIndex]
	player.Health -= amount
	if player.Health < 0 {
		player.Health = 0
	}

	return GameEvent{
		Type:     EventDamage,
		PlayerID: player.ID,
		Target: Target{
			Type:        TargetTypePlayer,
			Kind:        TargetKindHero,
			PlayerID:    player.ID,
			OwnerID:     player.ID,
			DisplayName: player.Name,
		},
		Amount:  amount,
		Message: "hero damaged by boss ability",
		Turn:    g.Turn,
	}
}

func damageMinion(g *Game, ownerIndex int, minionIndex int, amount int) GameEvent {
	if g == nil {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage to minion: game is nil",
		}
	}

	if ownerIndex < 0 || ownerIndex >= len(g.Players) {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage to minion: invalid owner index",
			Turn:    g.Turn,
		}
	}

	if minionIndex < 0 || minionIndex >= len(g.Players[ownerIndex].Board) {
		return GameEvent{
			Type:    EventDamage,
			Amount:  0,
			Message: "cannot deal damage to minion: invalid minion index",
			Turn:    g.Turn,
		}
	}

	if amount < 0 {
		amount = 0
	}

	minion := &g.Players[ownerIndex].Board[minionIndex]
	minion.Health -= amount
	if minion.Health < 0 {
		minion.Health = 0
	}

	return GameEvent{
		Type: EventDamage,
		Target: Target{
			Type:        TargetTypeMinion,
			MinionID:    minion.ID,
			OwnerID:     minion.OwnerID,
			DisplayName: minion.Name,
		},
		Amount:  amount,
		Message: "minion damaged by boss ability",
		Turn:    g.Turn,
	}
}
