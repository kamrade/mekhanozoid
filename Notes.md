newDevGame()
prepareDevGame()

type Action struct {
	Type     ActionType
	PlayerID PlayerID
	CardID   CardInstanceID
	SourceID MinionID
	TargetID string
	Target   Target
}

type runner {
  - main.go
    g: Game
      - ID - Game ID
      - Status           - GameStatus
      - Players          - []Player
      - Boss             - Boss
      - Turn             - Номер хода
      - ActivePlayerID   - ID активного игрока
      - Events           - []GameEvent
      - Seed             - int64

      - TargetIDHero0 = "hero:0"
	    - TargetIDHero1 = "hero:1"
    	- TargetIDBoss  = "boss"

      - ErrGameNotActive         = errors.New("game is not active")
      - ErrNotYourTurn           = errors.New("not your turn")
      - ErrUnknownAction         = errors.New("unknown action")
      - ErrCardNotInHand         = errors.New("card is not in hand")
      - ErrUnknownCard           = errors.New("unknown card")
      - ErrBoardFull             = errors.New("board is full")
      - ErrMinionNotFound        = errors.New("minion not found")
      - ErrMinionCantAttack      = errors.New("minion cannot attack")
      - ErrUnsupportedCardEffect = errors.New("unsupported card effect")
      - ErrTargetRequired        = errors.New("target is required")
      - ErrInvalidTarget         = errors.New("invalid target")

      - ActionTypeStartGame ActionType = "start_game"
      - ActionTypeEndTurn   ActionType = "end_turn"
      - ActionTypePlayCard  ActionType = "play_card"
      - ActionTypeAttack    ActionType = "attack"

      - ApplyAction()
      - applyEndTurn()
      - playCard()
      - applyAttack()

  - setup.go
    activePlayer()
    isGameOver()
    winnerName()

  - runner.go
    run()           - запускает обработку команд
    handleCommand() - обрабатвает команды
    cmdEnd()
    cmdPlay()
    cmdAttack()
    applyAndReport()

    cmdEnd()
    cmdPlay()
      - guards
      - 
    cmdAttack()

  - render.go
    renderState()
    renderHand()
    renderBoard()
    renderOtherBoard()
    renderRecentEvents()
    renderAllEvents()
    printHelp()
}

type Game {
  Status: GameStatusCreated | GameStatusActive | GameStatusWon | GameStatusLost
}