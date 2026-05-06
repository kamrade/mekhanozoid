package game

import "fmt"

const (
	StartingDeckSize = 10
	StartingHandSize = 3
)

var StartingCardDefinitions = []CardDefinition{
	{
		ID:          CardID("strike"),
		Name:        "Strike",
		Type:        CardTypeSpell,
		Cost:        1,
		Description: "Deal damage.",
	},
	{
		ID:          CardID("repair"),
		Name:        "Repair",
		Type:        CardTypeSpell,
		Cost:        1,
		Description: "Restore health.",
	},
	{
		ID:          CardID("guard_bot"),
		Name:        "Guard Bot",
		Type:        CardTypeMinion,
		Cost:        2,
		Description: "A small defensive minion.",
	},
	{
		ID:          CardID("blast"),
		Name:        "Blast",
		Type:        CardTypeSpell,
		Cost:        2,
		Description: "Deal more damage.",
	},
	{
		ID:          CardID("shield"),
		Name:        "Shield",
		Type:        CardTypeSpell,
		Cost:        1,
		Description: "Gain protection.",
	},
}

func NewStartingDeck(ownerID PlayerID) []CardInstance {
	deck := make([]CardInstance, 0, StartingDeckSize)

	for i := 0; i < StartingDeckSize; i++ {
		definition := StartingCardDefinitions[i%len(StartingCardDefinitions)]

		card := CardInstance{
			ID:           CardInstanceID(fmt.Sprintf("%s_card_%02d", ownerID, i+1)),
			DefinitionID: definition.ID,
			OwnerID:      ownerID,
		}

		deck = append(deck, card)
	}

	return deck
}

func drawStartingHand(deck []CardInstance, handSize int) ([]CardInstance, []CardInstance) {
	if handSize > len(deck) {
		handSize = len(deck)
	}

	hand := append([]CardInstance{}, deck[:handSize]...)
	remainingDeck := append([]CardInstance{}, deck[handSize:]...)

	return hand, remainingDeck
}
