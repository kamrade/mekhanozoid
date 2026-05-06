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
