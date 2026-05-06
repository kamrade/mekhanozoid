package game

import "errors"

var (
	ErrGameNotActive = errors.New("game is not active")
	ErrNotYourTurn   = errors.New("not your turn")
	ErrUnknownAction = errors.New("unknown action")
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
