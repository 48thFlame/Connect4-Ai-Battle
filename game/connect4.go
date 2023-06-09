package game

const (
	CRowsNum = 6
	CColsNum = 7
)

type CPlr int

const (
	CNone CPlr = iota
	CPlr1Max
	CPlr2Min
)

type CStateOfGame int

const (
	CStatePlaying CStateOfGame = iota
	CStatePlr1Won
	CStatePlr2Won
	CStateDraw
)

func NewConnect4Game() *Connect4Game {
	return &Connect4Game{
		Board:     CBoard{},
		PlrTurn:   CPlr1Max,
		GameState: CStatePlaying,
		TurnNum:   0,
	}
}

type CBoard [CRowsNum][CColsNum]CPlr
type Connect4Game struct {
	Board     CBoard       `json:"board"`   // an 2d array representing the board
	PlrTurn   CPlr         `json:"plrTurn"` // players turn 1/2
	GameState CStateOfGame `json:"gameState"`
	TurnNum   int          `json:"turnNumber"`
}

func (game *Connect4Game) Turn(col int) (good bool) {
	good = cMovePossible(game.Board, col)

	if good {
		game.Board = cPutPieceOnBoard(game.Board, col, game.PlrTurn)
		game.GameState = cGetGameState(game.Board)
		game.PlrTurn = cGetOtherPlayer(game.PlrTurn)
		game.TurnNum += 1
	}

	return good
}
