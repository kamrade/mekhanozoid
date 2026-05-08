package game

import (
	"strconv"
	"strings"
)

// ValidTargets returns all legal targets for a card instance in a player's hand.
// It is intended for future UI code that needs to highlight selectable targets.
func ValidTargets(g *Game, playerID string, cardInstanceID string) ([]Target, error) {
	if g == nil {
		return nil, ErrNilGame
	}

	playerIndex := findPlayerIndexByID(g, PlayerID(playerID))
	if playerIndex == -1 {
		return nil, ErrInvalidPlayerIndex
	}

	player := &g.Players[playerIndex]

	cardIndex := findCardInHand(player, CardInstanceID(cardInstanceID))
	if cardIndex == -1 {
		return nil, ErrCardNotInHand
	}

	cardInstance := player.Hand[cardIndex]

	cardDefinition, ok := CardRegistry[cardInstance.DefinitionID]
	if !ok {
		return nil, ErrUnknownCard
	}

	if !cardDefinition.Targeting.Required {
		return []Target{}, nil
	}

	return targetsForRule(g, cardDefinition.Targeting), nil
}

func ResolveTarget(g *Game, targetID string) (Target, int, error) {
	if g == nil {
		return Target{}, -1, ErrNilGame
	}

	switch targetID {
	case "":
		return Target{}, -1, ErrTargetRequired

	case TargetIDBoss:
		return Target{
			ID:          TargetIDBoss,
			Type:        TargetTypeBoss,
			Kind:        TargetKindBoss,
			BossID:      g.Boss.ID,
			DisplayName: g.Boss.Name,
		}, -1, nil

	case TargetIDHero0, TargetIDHero1:
		index, err := heroIndexFromTargetID(targetID)
		if err != nil {
			return Target{}, -1, ErrInvalidTarget
		}

		if index < 0 || index >= len(g.Players) {
			return Target{}, -1, ErrInvalidTarget
		}

		player := g.Players[index]

		return Target{
			ID:          targetID,
			Type:        TargetTypePlayer,
			Kind:        TargetKindHero,
			PlayerID:    player.ID,
			OwnerID:     player.ID,
			DisplayName: player.Name,
		}, index, nil

	default:
		minionTarget, ownerIndex, err := resolveMinionTarget(g, targetID)
		if err == nil {
			return minionTarget, ownerIndex, nil
		}

		return Target{}, -1, ErrInvalidTarget
	}
}

func ValidateTargetForCard(g *Game, card CardDefinition, targetID string) (Target, int, error) {
	if card.Targeting.Required && targetID == "" {
		return Target{}, -1, ErrTargetRequired
	}

	target, playerIndex, err := ResolveTarget(g, targetID)
	if err != nil {
		return Target{}, -1, err
	}

	if len(card.Targeting.AllowedKinds) == 0 {
		return target, playerIndex, nil
	}

	if !isAllowedTargetKind(target.Kind, card.Targeting.AllowedKinds) {
		return Target{}, -1, ErrInvalidTarget
	}

	return target, playerIndex, nil
}

func targetsForRule(g *Game, rule TargetingRule) []Target {
	if g == nil {
		return []Target{}
	}

	targets := []Target{}

	for _, kind := range rule.AllowedKinds {
		switch kind {
		case TargetKindHero:
			targets = append(targets, heroTargets(g)...)

		case TargetKindBoss:
			targets = append(targets, bossTarget(g))

		case TargetKindMinion:
			targets = append(targets, livingMinionTargets(g)...)
		}
	}

	return targets
}

func heroTargets(g *Game) []Target {
	targets := make([]Target, 0, len(g.Players))

	for i, player := range g.Players {
		targetID := "hero:" + strconv.Itoa(i)

		targets = append(targets, Target{
			ID:          targetID,
			Type:        TargetTypePlayer,
			Kind:        TargetKindHero,
			PlayerID:    player.ID,
			OwnerID:     player.ID,
			DisplayName: player.Name,
		})
	}

	return targets
}

func bossTarget(g *Game) Target {
	return Target{
		ID:          TargetIDBoss,
		Type:        TargetTypeBoss,
		Kind:        TargetKindBoss,
		BossID:      g.Boss.ID,
		DisplayName: g.Boss.Name,
	}
}

func isAllowedTargetKind(kind TargetKind, allowed []TargetKind) bool {
	for _, allowedKind := range allowed {
		if kind == allowedKind {
			return true
		}
	}

	return false
}

func heroIndexFromTargetID(targetID string) (int, error) {
	const prefix = "hero:"

	if len(targetID) <= len(prefix) || targetID[:len(prefix)] != prefix {
		return -1, ErrInvalidTarget
	}

	indexText := targetID[len(prefix):]

	index, err := strconv.Atoi(indexText)
	if err != nil {
		return -1, ErrInvalidTarget
	}

	return index, nil
}

func livingMinionTargets(g *Game) []Target {
	if g == nil {
		return nil
	}

	targets := []Target{}

	for ownerIndex := range g.Players {
		player := g.Players[ownerIndex]
		for minionIndex := range player.Board {
			minion := player.Board[minionIndex]
			if minion.Health <= 0 {
				continue
			}

			targets = append(targets, Target{
				ID:          minionTargetID(ownerIndex, minion.ID),
				Type:        TargetTypeMinion,
				Kind:        TargetKindMinion,
				MinionID:    minion.ID,
				OwnerID:     minion.OwnerID,
				DisplayName: minion.Name,
			})
		}
	}

	return targets
}

func resolveMinionTarget(g *Game, targetID string) (Target, int, error) {
	const prefix = "minion:"

	if g == nil {
		return Target{}, -1, ErrNilGame
	}

	if !strings.HasPrefix(targetID, prefix) {
		return Target{}, -1, ErrInvalidTarget
	}

	parts := strings.Split(targetID, ":")
	if len(parts) != 3 {
		return Target{}, -1, ErrInvalidTarget
	}

	ownerIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return Target{}, -1, ErrInvalidTarget
	}

	if ownerIndex < 0 || ownerIndex >= len(g.Players) {
		return Target{}, -1, ErrInvalidTarget
	}

	minionID := MinionID(parts[2])
	if minionID == "" {
		return Target{}, -1, ErrInvalidTarget
	}

	for i := range g.Players[ownerIndex].Board {
		minion := g.Players[ownerIndex].Board[i]
		if minion.ID != minionID {
			continue
		}

		if minion.Health <= 0 {
			return Target{}, -1, ErrInvalidTarget
		}

		return Target{
			ID:          targetID,
			Type:        TargetTypeMinion,
			Kind:        TargetKindMinion,
			MinionID:    minion.ID,
			OwnerID:     minion.OwnerID,
			DisplayName: minion.Name,
		}, ownerIndex, nil
	}

	return Target{}, -1, ErrInvalidTarget
}
