// This file contains the card registry.
// The registry is the source of truth for all card definitions available to the game engine.

package game

var CardRegistry = map[CardID]CardDefinition{
	CardID("strike"): {
		ID:          CardID("strike"),
		Name:        "Strike",
		Type:        CardTypeSpell,
		Cost:        1,
		Description: "Deal 3 damage to the boss.",
		Effect: CardEffect{
			Type:   EffectDamageBoss,
			Amount: 3,
		},
		Targeting: TargetingRule{
			Required: false,
			AllowedKinds: []TargetKind{
				TargetKindBoss,
			},
		},
	},
	CardID("repair"): {
		ID:          CardID("repair"),
		Name:        "Repair",
		Type:        CardTypeSpell,
		Cost:        2,
		Description: "Restore 5 health to a hero.",
		Effect: CardEffect{
			Type:   EffectHealHero,
			Amount: 5,
		},
		Targeting: TargetingRule{
			Required: true,
			AllowedKinds: []TargetKind{
				TargetKindHero,
			},
		},
	},
	CardID("drone"): {
		ID:          CardID("drone"),
		Name:        "Drone",
		Type:        CardTypeMinion,
		Cost:        2,
		Description: "Summon a 2/3 Drone.",
		Effect: CardEffect{
			Type: EffectSummon,
		},
		Attack: 2,
		Health: 3,
	},
}
