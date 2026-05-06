package game

import "strconv"

func ResolveTarget(g *Game, targetID string) (Target, int, error) {
	if g == nil {
		return Target{}, -1, ErrNilGame
	}

	switch targetID {
	case "":
		return Target{}, -1, ErrTargetRequired

	case TargetIDBoss:
		return Target{
			Type:   TargetTypeBoss,
			BossID: g.Boss.ID,
		}, -1, nil

	case TargetIDHero0, TargetIDHero1:
		index, err := heroIndexFromTargetID(targetID)
		if err != nil {
			return Target{}, -1, ErrInvalidTarget
		}

		if index < 0 || index >= len(g.Players) {
			return Target{}, -1, ErrInvalidTarget
		}

		return Target{
			Type:     TargetTypePlayer,
			PlayerID: g.Players[index].ID,
		}, index, nil

	default:
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

	targetKind := targetKindFromTarget(target)
	if !isAllowedTargetKind(targetKind, card.Targeting.AllowedKinds) {
		return Target{}, -1, ErrInvalidTarget
	}

	return target, playerIndex, nil
}

func targetKindFromTarget(target Target) TargetKind {
	switch target.Type {
	case TargetTypePlayer:
		return TargetKindHero
	case TargetTypeBoss:
		return TargetKindBoss
	default:
		return TargetKindNone
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
