// ids.go
// This file defines strongly typed identifiers used by the game domain.
// Typed IDs help avoid mixing unrelated entities such as players, cards, bosses, and minions.

// game.go
// This file defines the root Game aggregate and game lifecycle statuses.
// It also provides NewGame, the factory function for creating the initial game state.

// player_config.go
// This file defines PlayerConfig, the input data required to create players
// when initializing a new game.

// player.go
// This file defines the Player domain model.
// A player owns cards, has health and mana, and may be marked as the current active player.

// boss.go
// This file defines the Boss domain model.
// The boss represents the shared enemy that both players cooperate against.

// card.go
// This file defines card-related domain types.
// CardDefinition describes a card template, while CardInstance represents a concrete card in a game.

// card_registry.go
// This file contains the card registry.
// The registry is the source of truth for all card definitions available to the game engine.

// card_validation.go
// This file contains validation helpers for card definitions and card instances.
// These helpers ensure that decks and hands only reference cards known by the registry.

// deck.go
// This file contains deck-related logic.
// It defines the starting deck, creates card instances, deals starting hands, and supports drawing cards.

// shuffle.go
// This file contains deterministic shuffle logic for card slices.
// The shuffle uses a seed so game setup can be reproduced in tests and debugging.

// minion.go
// This file defines the Minion domain model.
// A minion is a summoned board entity owned by a player and created from a card definition.

// action.go
// This file defines player actions and action targets.
// Actions describe player intent and are processed by the game engine through ApplyAction.

// event.go
// This file defines game events.
// Events describe important changes in game state and can later be used by UI, logs, or replays.

// mana.go
// This file contains mana-related rules.
// It supports refreshing mana at the start of a turn and spending mana safely.

// engine.go
// This file contains the central game engine entry point.
// ApplyAction validates and applies player actions, mutates game state, and returns generated events.

// This file defines the root Game aggregate and game lifecycle statuses.
// It also provides NewGame, the factory function for creating the initial game state.

package game

type GameStatus string

const (
	GameStatusCreated GameStatus = "created"
	GameStatusActive  GameStatus = "active"
	GameStatusWon     GameStatus = "won"
	GameStatusLost    GameStatus = "lost"
)

type Game struct {
	ID             GameID
	Status         GameStatus
	Players        []Player
	Boss           Boss
	Turn           int
	ActivePlayerID PlayerID
	Events         []GameEvent
	Seed           int64
}

func NewGame(id string, p1 PlayerConfig, p2 PlayerConfig, seed int64) *Game {
	player1 := NewPlayer(PlayerID(p1.ID), p1.Name)
	player2 := NewPlayer(PlayerID(p2.ID), p2.Name)

	player1.Deck = NewStartingDeck(player1.ID, seed)
	player2.Deck = NewStartingDeck(player2.ID, seed+1)

	ShuffleCards(player1.Deck, seed)
	ShuffleCards(player2.Deck, seed+1)

	player1.Hand, player1.Deck = drawStartingHand(player1.Deck, StartingHandSize)
	player2.Hand, player2.Deck = drawStartingHand(player2.Deck, StartingHandSize)

	player1.IsCurrent = true
	player2.IsCurrent = false

	return &Game{
		ID:             GameID(id),
		Status:         GameStatusActive,
		Players:        []Player{player1, player2},
		Boss:           NewBoss(BossID("boss_1"), "Mechazod"),
		Turn:           1,
		ActivePlayerID: player1.ID,
		Events:         []GameEvent{},
		Seed:           seed,
	}
}
