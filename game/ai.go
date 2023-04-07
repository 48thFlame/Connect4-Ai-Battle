package game

import (
	"fmt"
	"time"
)

func FlameAiGetMove(game Connect4Game) (col int) {

	fStats = &flameStats{}

	s := make([]moveRating, 0)
	moves := getOrderedAvailableMoves(game.Board)
	depth := 2
	timeout := time.After(timeoutTime)

loop:
	for {
		select {
		case <-timeout:
			break loop
		default:
		}
		depth += 3
		s = search(game, depth, moves)
		moves = moveRatingsToMoves(s)
	}

	fmt.Printf("Max-search at depth %v\n", depth)
	bestMove := s[0]

	col = bestMove.move
	eval := bestMove.eval

	fmt.Printf("%v\n", s)
	fmt.Printf(
		"Eval: %v | Analyzed: %v | Broke: %v \n",
		eval,
		fStats.posAnalyzed,
		fStats.posBroke)

	return col
}

func search(game Connect4Game, depth int, moves []int) (s []moveRating) {

	isMaximizingPlayer := game.PlrTurn == CPlr1Max
	evalCh := make(chan moveRating)

	for _, moveCol := range moves {
		go func(m int) {
			// opposite of max because is placing now the other one
			posEval := minMax(cPutPieceOnBoard(game.Board, m, game.PlrTurn), depth, -highNumber, highNumber, !isMaximizingPlayer)
			evalCh <- moveRating{move: m, eval: posEval}
		}(moveCol)
	}

	numOfPossible := len(moves)
	for i := 0; i < numOfPossible; i++ { // for each goroutine started, get its result
		t := <-evalCh

		s = append(s, t)
	}

	if isMaximizingPlayer {
		sortMovesHighLow(s)
	} else {
		sortMovesLowHigh(s)
	}

	return s
}
