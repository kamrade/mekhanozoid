package game

import "fmt"

const (
	StartingDeckSize = 10
	StartingHandSize = 3
)

var StartingDeckCardIDs = []CardID{
	CardID("strike"),
	CardID("repair"),
	CardID("drone"),
	CardID("strike"),
	CardID("repair"),
	CardID("drone"),
	CardID("strike"),
	CardID("repair"),
	CardID("drone"),
	CardID("strike"),
}

func NewStartingDeck(ownerID PlayerID) []CardInstance {
	deck := make([]CardInstance, 0, len(StartingDeckCardIDs))

	for i, cardID := range StartingDeckCardIDs {
		card := CardInstance{
			ID:           CardInstanceID(fmt.Sprintf("%s_card_%02d", ownerID, i+1)),
			DefinitionID: cardID,
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
