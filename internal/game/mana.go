package game

import "errors"

const (
	StartingMana    = 0
	StartingMaxMana = 0
	MaxMana         = 10
)

var (
	ErrInvalidPlayerIndex = errors.New("invalid player index")
	ErrNotEnoughMana      = errors.New("not enough mana")
	ErrNegativeManaAmount = errors.New("mana amount cannot be negative")
	ErrNilGame            = errors.New("game is nil")
)

func RefreshMana(g *Game, playerIndex int) {
	if g == nil {
		return
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return
	}

	player := &g.Players[playerIndex]

	if player.MaxMana < MaxMana {
		player.MaxMana++
	}

	player.Mana = player.MaxMana
}

func SpendMana(g *Game, playerIndex int, amount int) error {
	if g == nil {
		return ErrNilGame
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return ErrInvalidPlayerIndex
	}

	if amount < 0 {
		return ErrNegativeManaAmount
	}

	player := &g.Players[playerIndex]

	if amount > player.Mana {
		return ErrNotEnoughMana
	}

	player.Mana -= amount

	return nil
}
