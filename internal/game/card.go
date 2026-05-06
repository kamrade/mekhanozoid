package game

type CardType string

const (
	CardTypeSpell  CardType = "spell"
	CardTypeMinion CardType = "minion"
	CardTypeWeapon CardType = "weapon"
)

type CardDefinition struct {
	ID          CardID
	Name        string
	Type        CardType
	Cost        int
	Description string
}

type CardInstance struct {
	ID           CardInstanceID
	DefinitionID CardID
	OwnerID      PlayerID
}
