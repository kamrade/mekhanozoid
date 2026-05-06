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

func DrawCard(g *Game, playerIndex int) []GameEvent {
	if g == nil {
		return []GameEvent{
			{
				Type:    EventTypeCardDrawn,
				Message: "cannot draw card: game is nil",
			},
		}
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return []GameEvent{
			{
				Type:    EventTypeCardDrawn,
				Message: "cannot draw card: invalid player index",
				Turn:    g.Turn,
			},
		}
	}

	player := &g.Players[playerIndex]

	if len(player.Deck) == 0 {
		return []GameEvent{
			{
				Type:     EventTypeCardDrawn,
				PlayerID: player.ID,
				Message:  "cannot draw card: deck is empty",
				Turn:     g.Turn,
			},
		}
	}

	card := player.Deck[0]

	player.Deck = player.Deck[1:]
	player.Hand = append(player.Hand, card)

	event := GameEvent{
		Type:     EventTypeCardDrawn,
		PlayerID: player.ID,
		CardID:   card.ID,
		Message:  "card drawn",
		Turn:     g.Turn,
	}

	g.Events = append(g.Events, event)

	return []GameEvent{event}
}
