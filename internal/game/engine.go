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
	ErrBoardFull             = errors.New("board is full")
	ErrMinionNotFound        = errors.New("minion not found")
	ErrMinionCantAttack      = errors.New("minion cannot attack")
	ErrUnsupportedCardEffect = errors.New("unsupported card effect")
	ErrTargetRequired        = errors.New("target is required")
	ErrInvalidTarget         = errors.New("invalid target")
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
	case ActionTypeAttack:
		return applyAttack(g, action)
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
	RefreshMinions(g, newActivePlayerIndex)
	bossAbilityEvents := ResolveBossAbility(g)
	g.Events = append(g.Events, bossAbilityEvents...)
	gameOverEvents := CheckGameOver(g)

	if len(gameOverEvents) > 0 {
		events := make([]GameEvent, 0, len(bossAbilityEvents)+len(gameOverEvents))
		events = append(events, bossAbilityEvents...)
		events = append(events, gameOverEvents...)
		return events, nil
	}

	turnStartedEvent := GameEvent{
		Type:     EventTypeTurnStarted,
		PlayerID: g.ActivePlayerID,
		Message:  "turn started",
		Turn:     g.Turn,
	}

	g.Events = append(g.Events, turnStartedEvent)

	drawEvents := DrawCard(g, newActivePlayerIndex)

	events := make([]GameEvent, 0, len(bossAbilityEvents)+1+len(drawEvents))
	events = append(events, bossAbilityEvents...)
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

	if player.Mana < cardDefinition.Cost {
		return nil, ErrNotEnoughMana
	}

	var resolvedTarget Target
	var targetPlayerIndex int

	switch cardDefinition.Type {
	case CardTypeSpell:
		var err error
		if cardDefinition.Targeting.Required || action.TargetID != "" {
			resolvedTarget, targetPlayerIndex, err = ValidateTargetForCard(g, cardDefinition, action.TargetID)
			if err != nil {
				return nil, err
			}
		}

		if err := validateSupportedSpellEffect(cardDefinition); err != nil {
			return nil, err
		}
	case CardTypeMinion:
		if err := canSummonMinion(player); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedCardEffect
	}

	if err := SpendMana(g, playerIndex, cardDefinition.Cost); err != nil {
		return nil, err
	}

	player.Hand = removeCardFromHand(player.Hand, cardIndex)

	cardPlayedEvent := GameEvent{
		Type:     EventCardPlayed,
		PlayerID: player.ID,
		CardID:   cardInstance.ID,
		Target:   resolvedTarget,
		Message:  "card played",
		Turn:     g.Turn,
	}

	events := []GameEvent{
		cardPlayedEvent,
	}

	switch cardDefinition.Type {
	case CardTypeSpell:
		switch cardDefinition.Effect.Type {
		case EffectDamageBoss:
			damageEvent := DealDamage(g, cardDefinition.Effect.Amount)
			damageEvent.PlayerID = player.ID
			damageEvent.CardID = cardInstance.ID

			events = append(events, damageEvent)
		case EffectHealHero:
			healEvent := Heal(g, targetPlayerIndex, cardDefinition.Effect.Amount)
			healEvent.CardID = cardInstance.ID

			events = append(events, healEvent)
		default:
			return nil, ErrUnsupportedCardEffect
		}
	case CardTypeMinion:
		summonEvent, err := SummonMinion(g, playerIndex, cardDefinition, cardInstance)
		if err != nil {
			return nil, err
		}

		events = append(events, summonEvent)
	}

	g.Events = append(g.Events, events...)

	gameOverEvents := CheckGameOver(g)
	events = append(events, gameOverEvents...)

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

func validateSupportedSpellEffect(card CardDefinition) error {
	switch card.Effect.Type {
	case EffectDamageBoss, EffectHealHero:
		return nil
	default:
		return ErrUnsupportedCardEffect
	}
}

func applyAttack(g *Game, action Action) ([]GameEvent, error) {
	playerIndex := findPlayerIndexByID(g, action.PlayerID)
	if playerIndex == -1 {
		return nil, ErrInvalidPlayerIndex
	}

	if g.ActivePlayerID != action.PlayerID {
		return nil, ErrNotYourTurn
	}

	if action.TargetID != TargetIDBoss {
		return nil, ErrInvalidTarget
	}

	minionIndex := findMinionOnBoard(&g.Players[playerIndex], action.SourceID)
	if minionIndex == -1 {
		return nil, ErrMinionNotFound
	}

	minion := &g.Players[playerIndex].Board[minionIndex]
	if !minion.CanAttack {
		return nil, ErrMinionCantAttack
	}

	minion.CanAttack = false

	attackEvent := GameEvent{
		Type:     EventAttack,
		PlayerID: action.PlayerID,
		SourceID: minion.ID,
		Target: Target{
			ID:          TargetIDBoss,
			Type:        TargetTypeBoss,
			Kind:        TargetKindBoss,
			BossID:      g.Boss.ID,
			DisplayName: g.Boss.Name,
		},
		Amount:  minion.Attack,
		Message: "minion attacked",
		Turn:    g.Turn,
	}

	damageEvent := DealDamage(g, minion.Attack)
	damageEvent.PlayerID = action.PlayerID
	damageEvent.SourceID = minion.ID

	events := []GameEvent{attackEvent, damageEvent}
	g.Events = append(g.Events, events...)

	gameOverEvents := CheckGameOver(g)
	events = append(events, gameOverEvents...)

	return events, nil
}

func findMinionOnBoard(player *Player, minionID MinionID) int {
	if player == nil || minionID == "" {
		return -1
	}

	for i := range player.Board {
		if player.Board[i].ID == minionID {
			return i
		}
	}

	return -1
}

func RefreshMinions(g *Game, playerIndex int) {
	if g == nil {
		return
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return
	}

	for i := range g.Players[playerIndex].Board {
		g.Players[playerIndex].Board[i].CanAttack = true
	}
}
