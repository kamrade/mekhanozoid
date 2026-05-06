package game

type GameStatus string

const (
	GameStatusCreated  GameStatus = "created"
	GameStatusRunning  GameStatus = "running"
	GameStatusWon      GameStatus = "won"
	GameStatusLost     GameStatus = "lost"
	GameStatusFinished GameStatus = "finished"
)

type Game struct {
	ID      GameID
	Status  GameStatus
	Players []Player
	Boss    Boss
	Turn    int
	Events  []GameEvent
}

func NewGame() *Game {
	return &Game{
		ID:     GameID("game_1"),
		Status: GameStatusCreated,
		Players: []Player{
			NewPlayer(PlayerID("player_1"), "Player 1"),
			NewPlayer(PlayerID("player_2"), "Player 2"),
		},
		Boss:   NewBoss(BossID("boss_1"), "Mechazod"),
		Turn:   0,
		Events: []GameEvent{},
	}
}
