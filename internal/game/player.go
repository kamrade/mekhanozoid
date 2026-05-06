package game

type Player struct {
	ID        PlayerID
	Name      string
	Health    int
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
		Health:    30,
		Deck:      []CardInstance{},
		Hand:      []CardInstance{},
		Board:     []Minion{},
		Discard:   []CardInstance{},
		IsCurrent: false,
	}
}
