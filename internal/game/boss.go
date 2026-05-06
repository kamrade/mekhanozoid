package game

const StartingBossHealth = 95

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
		Health:    StartingBossHealth,
		Attack:    2,
		Armor:     0,
		MaxHealth: StartingBossHealth,
	}
}
