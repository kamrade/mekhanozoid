// This file defines card-related domain types.
// CardDefinition describes a card template, while CardInstance represents a concrete card in a game.

package game

type CardType string

const (
	CardTypeSpell  CardType = "spell"
	CardTypeMinion CardType = "minion"
)

type EffectType string

const (
	EffectNone       EffectType = "none"
	EffectDamageBoss EffectType = "damage_boss"
	EffectHealHero   EffectType = "heal_hero"
	EffectSummon     EffectType = "summon"
)

type CardEffect struct {
	Type   EffectType
	Amount int
}

type CardDefinition struct {
	ID          CardID
	Name        string
	Type        CardType
	Cost        int
	Description string

	Effect CardEffect

	Attack int
	Health int
}

type CardInstance struct {
	ID           CardInstanceID
	DefinitionID CardID
	OwnerID      PlayerID
}
