package game

import (
	"fmt"
	"math"
)

func FlameAiGetMove(game Connect4Game) (col int) {
	possibleMoves := getOrderedMoves(game.Board)

	fStats = &flameStats{}

	var depth int

	if game.TurnNum < 6 {
		depth = 7
	} else if game.TurnNum < 12 {
		depth = 8
	} else if game.TurnNum < 20 {
		depth = 9
	} else if game.TurnNum < 30 {
		depth = 12
	} else {
		depth = game.TurnNum
	}

	isMaximizingPlayer := game.PlrTurn == CPlr1Max
	evalCh := make(chan [2]int)

	for _, moveCol := range possibleMoves {
		go func(m int) {
			// opposite of max because is placing now the other one
			posEval := minMax(cPutPieceOnBoard(game.Board, m, game.PlrTurn), depth, -highNumber, highNumber, !isMaximizingPlayer)
			evalCh <- [2]int{m, posEval}
		}(moveCol)
	}

	var extremeEval int
	if isMaximizingPlayer {
		extremeEval = -math.MaxInt
	} else {
		extremeEval = math.MaxInt
	}
	var bestMove int

	numOfPossible := len(possibleMoves)
	for i := 0; i < numOfPossible; i++ { // for each goroutine started, get its result
		t := <-evalCh

		posEval := t[1]

		if isMaximizingPlayer {
			if posEval >= extremeEval {
				extremeEval = posEval
				bestMove = t[0]
			}
		} else {
			if posEval < extremeEval {
				extremeEval = posEval
				bestMove = t[0]
			}
		}
	}

	fmt.Printf(
		"Eval: %v | Analyzed: %v | Broke: %v \n",
		extremeEval,
		fStats.posAnalyzed,
		fStats.posBroke)

	return bestMove
}
