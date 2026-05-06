package game

type Status string

const (
	StatusCreated Status = "created"
)

type Game struct {
	status Status
}

func NewGame() *Game {
	return &Game{
		status: StatusCreated,
	}
}

func (g *Game) Status() Status {
	return g.status
}
