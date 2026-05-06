// This file defines the Minion domain model.
// A minion is a summoned board entity owned by a player and created from a card definition.

package game

type Minion struct {
	ID           MinionID
	DefinitionID CardID
	OwnerID      PlayerID
	Name         string
	Attack       int
	Health       int
	MaxHealth    int
	Exhausted    bool
}
