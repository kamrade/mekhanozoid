// This file contains deck-related logic.
// It defines the starting deck, creates card instances, deals starting hands, and supports drawing cards.

package game

import (
	"fmt"
	"math/rand"
)

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

func NewStartingDeck(ownerID PlayerID, seed int64) []CardInstance {
	deck := make([]CardInstance, 0, len(StartingDeckCardIDs))
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < StartingDeckSize; i++ {
		cardID := StartingDeckCardIDs[r.Intn(len(StartingDeckCardIDs))]
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
		player.FatigueDamage++
		fatigue := player.FatigueDamage

		player.Health -= fatigue
		if player.Health < 0 {
			player.Health = 0
		}

		drawEvent := GameEvent{
			Type:     EventTypeCardDrawn,
			PlayerID: player.ID,
			Message:  "cannot draw card: deck is empty",
			Turn:     g.Turn,
		}
		damageEvent := GameEvent{
			Type:     EventDamage,
			PlayerID: player.ID,
			Target: Target{
				Type:        TargetTypePlayer,
				Kind:        TargetKindHero,
				PlayerID:    player.ID,
				OwnerID:     player.ID,
				DisplayName: player.Name,
			},
			Amount:  fatigue,
			Message: "fatigue damage: deck is empty",
			Turn:    g.Turn,
		}

		g.Events = append(g.Events, drawEvent, damageEvent)
		return []GameEvent{drawEvent, damageEvent}
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
