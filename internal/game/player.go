// This file defines the Player domain model.
// A player owns cards, has health and mana, and may be marked as the current active player.

package game

const StartingPlayerHealth = 30

type Player struct {
	ID        PlayerID
	Name      string
	Health    int
	MaxHealth int
	FatigueDamage int
	Mana      int
	MaxMana   int
	Deck      []CardInstance
	Hand      []CardInstance
	Board     []Minion
	Discard   []CardInstance
	IsCurrent bool
}

func NewPlayer(id PlayerID, name string) Player {
	return Player{
		ID:        id,
		Name:      name,
		Health:    StartingPlayerHealth,
		MaxHealth: StartingPlayerHealth,
		FatigueDamage: 0,
		Mana:      StartingMana,
		MaxMana:   StartingMaxMana,
		Deck:      []CardInstance{},
		Hand:      []CardInstance{},
		Board:     []Minion{},
		Discard:   []CardInstance{},
		IsCurrent: false,
	}
}
