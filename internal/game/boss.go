package game

type Boss struct {
	ID        BossID
	Name      string
	Health    int
	Attack    int
	Armor     int
	MaxHealth int
}

func NewBoss(id BossID, name string) Boss {
	return Boss{
		ID:        id,
		Name:      name,
		Health:    95,
		Attack:    2,
		Armor:     0,
		MaxHealth: 95,
	}
}
