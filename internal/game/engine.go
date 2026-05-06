// This file contains the central game engine entry point.
// ApplyAction validates and applies player actions, mutates game state, and returns generated events.

package game

import "errors"

var (
	ErrGameNotActive         = errors.New("game is not active")
	ErrNotYourTurn           = errors.New("not your turn")
	ErrUnknownAction         = errors.New("unknown action")
	ErrCardNotInHand         = errors.New("card is not in hand")
	ErrUnknownCard           = errors.New("unknown card")
	ErrUnsupportedCardEffect = errors.New("unsupported card effect")
)

func ApplyAction(g *Game, action Action) ([]GameEvent, error) {
	if g == nil {
		return nil, ErrNilGame
	}

	if g.Status != GameStatusActive {
		return nil, ErrGameNotActive
	}

	switch action.Type {
	case ActionTypeEndTurn:
		return applyEndTurn(g, action)
	case ActionTypePlayCard:
		return playCard(g, action)
	default:
		return nil, ErrUnknownAction
	}
}

func applyEndTurn(g *Game, action Action) ([]GameEvent, error) {
	currentPlayerIndex := findPlayerIndexByID(g, action.PlayerID)
	if currentPlayerIndex == -1 {
		return nil, ErrInvalidPlayerIndex
	}

	if g.ActivePlayerID != action.PlayerID {
		return nil, ErrNotYourTurn
	}

	newActivePlayerIndex := nextPlayerIndex(currentPlayerIndex, len(g.Players))

	setActivePlayer(g, newActivePlayerIndex)

	g.Turn++

	RefreshMana(g, newActivePlayerIndex)

	turnStartedEvent := GameEvent{
		Type:     EventTypeTurnStarted,
		PlayerID: g.ActivePlayerID,
		Message:  "turn started",
		Turn:     g.Turn,
	}

	g.Events = append(g.Events, turnStartedEvent)

	drawEvents := DrawCard(g, newActivePlayerIndex)

	events := make([]GameEvent, 0, 1+len(drawEvents))
	events = append(events, turnStartedEvent)
	events = append(events, drawEvents...)

	return events, nil
}

func playCard(g *Game, action Action) ([]GameEvent, error) {
	playerIndex := findPlayerIndexByID(g, action.PlayerID)
	if playerIndex == -1 {
		return nil, ErrInvalidPlayerIndex
	}

	if g.ActivePlayerID != action.PlayerID {
		return nil, ErrNotYourTurn
	}

	player := &g.Players[playerIndex]

	cardIndex := findCardInHand(player, action.CardID)
	if cardIndex == -1 {
		return nil, ErrCardNotInHand
	}

	cardInstance := player.Hand[cardIndex]

	cardDefinition, ok := CardRegistry[cardInstance.DefinitionID]
	if !ok {
		return nil, ErrUnknownCard
	}

	if cardDefinition.Type != CardTypeSpell {
		return nil, ErrUnsupportedCardEffect
	}

	if cardDefinition.Effect.Type != EffectDamageBoss {
		return nil, ErrUnsupportedCardEffect
	}

	if player.Mana < cardDefinition.Cost {
		return nil, ErrNotEnoughMana
	}

	if err := SpendMana(g, playerIndex, cardDefinition.Cost); err != nil {
		return nil, err
	}

	player.Hand = removeCardFromHand(player.Hand, cardIndex)

	cardPlayedEvent := GameEvent{
		Type:     EventCardPlayed,
		PlayerID: player.ID,
		CardID:   cardInstance.ID,
		Message:  "card played",
		Turn:     g.Turn,
	}

	damageEvent := DealDamage(g, cardDefinition.Effect.Amount)
	damageEvent.PlayerID = player.ID
	damageEvent.CardID = cardInstance.ID

	events := []GameEvent{
		cardPlayedEvent,
		damageEvent,
	}

	g.Events = append(g.Events, events...)

	return events, nil
}

func findCardInHand(player *Player, cardID CardInstanceID) int {
	if player == nil {
		return -1
	}

	for i, card := range player.Hand {
		if card.ID == cardID {
			return i
		}
	}

	return -1
}

func removeCardFromHand(hand []CardInstance, index int) []CardInstance {
	if index < 0 || index >= len(hand) {
		return hand
	}

	result := make([]CardInstance, 0, len(hand)-1)
	result = append(result, hand[:index]...)
	result = append(result, hand[index+1:]...)

	return result
}

func findPlayerIndexByID(g *Game, playerID PlayerID) int {
	if g == nil {
		return -1
	}

	for i, player := range g.Players {
		if player.ID == playerID {
			return i
		}
	}

	return -1
}

func nextPlayerIndex(currentIndex int, playerCount int) int {
	if playerCount == 0 {
		return -1
	}

	return (currentIndex + 1) % playerCount
}

func setActivePlayer(g *Game, playerIndex int) {
	if g == nil {
		return
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return
	}

	for i := range g.Players {
		g.Players[i].IsCurrent = i == playerIndex
	}

	g.ActivePlayerID = g.Players[playerIndex].ID
}
