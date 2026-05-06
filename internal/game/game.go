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

	player1.Deck = NewStartingDeck(player1.ID)
	player2.Deck = NewStartingDeck(player2.ID)

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
