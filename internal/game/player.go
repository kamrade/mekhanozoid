package game

const StartingPlayerHealth = 30

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
		Health:    StartingPlayerHealth,
		Deck:      []CardInstance{},
		Hand:      []CardInstance{},
		Board:     []Minion{},
		Discard:   []CardInstance{},
		IsCurrent: false,
	}
}
