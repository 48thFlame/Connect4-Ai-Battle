package game

import (
	"fmt"
	"time"
)

func FlameAiGetMove(game Connect4Game) (col int) {
	fStats = &flameStats{}

	mRatings := make([]moveRating, 0)
	moves := getOrderedAvailableMoves(game.Board)
	depth := 6
	timeout := time.After(timeoutTime)

loop:
	for {
		select {
		case <-timeout:
			break loop
		default:
		}
		depth += 1
		mRatings = search(game, depth, moves)
		moves = moveRatingsToMoves(mRatings)
	}

	var eval int
	if len(mRatings) > 0 { //actually made a search (almost definitely unless timed-out first)
		fmt.Printf("Max-search at depth %v\n", depth)
		bestMove := mRatings[0]

		eval = bestMove.eval
	}

	fmt.Printf("%v\n", mRatings)
	fmt.Printf(
		"Eval: %v | Analyzed: %v | Broke: %v \n",
		eval,
		fStats.posAnalyzed,
		fStats.posBroke)

	return moves[0]
}

func search(game Connect4Game, depth int, moves []int) (s []moveRating) {
	isMaximizingPlayer := game.PlrTurn == CPlr1Max
	evalCh := make(chan moveRating)

	for _, moveCol := range moves {
		go func(m int) {
			// opposite of max because is placing now the other one
			posEval := minMax(
				cPutPieceOnBoard(game.Board, m, game.PlrTurn),
				depth, -highNumber, highNumber, !isMaximizingPlayer)

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

func minMax(board CBoard, depth int, alpha int, beta int, maximizingPlayer bool) int {
	if gs := cGetGameState(board); gs != CStatePlaying {
		fStats.posAnalyzed++

		switch gs {
		case CStateDraw:
			return 0
		case CStatePlr1Won:
			return highNumber - depth
		case CStatePlr2Won:
			return -highNumber - depth
		}
	} else if depth == 0 {
		fStats.posAnalyzed++

		return staticEval(&board)
	}

	possibleMoves := getOrderedAvailableMoves(board)
	if maximizingPlayer {
		maxEval := -highNumber

		for _, move := range possibleMoves {
			newEval := minMax(cPutPieceOnBoard(board, move, CPlr1Max), depth-1, alpha, beta, false)
			maxEval = max(maxEval, newEval)

			alpha = max(alpha, newEval)
			if beta <= alpha {
				fStats.posBroke++
				break
			}
		}

		return maxEval
	} else {
		minEval := highNumber

		for _, move := range possibleMoves {
			newEval := minMax(cPutPieceOnBoard(board, move, CPlr2Min), depth-1, alpha, beta, true)
			minEval = min(minEval, newEval)

			beta = min(beta, newEval)
			if beta <= alpha {
				fStats.posBroke++
				break
			}
		}

		return minEval
	}
}
