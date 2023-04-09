package game

import (
	"fmt"
)

func FlameAiGetMove(game Connect4Game) (col int) {
	var depth int

	if game.TurnNum < 3 {
		depth = 7
	} else if game.TurnNum < 14 {
		depth = 8
	} else if game.TurnNum < 22 {
		depth = 9
	} else if game.TurnNum < 26 {
		depth = 11
	} else {
		depth = game.TurnNum * game.TurnNum
	}

	initOrdered := getOrderedAvailableMoves(game.Board)
	quickSearchResults := moveRatingsToMoves(search(game, 4, initOrdered))
	fStats = &flameStats{}
	mRatings := search(game, depth, quickSearchResults)

	bestMove := mRatings[0]

	eval := bestMove.eval

	fmt.Printf("%v\n", mRatings)
	fmt.Printf(
		"Depth: %v | Eval: %v | Analyzed: %v | Broke: %v \n",
		depth,
		eval,
		fStats.posAnalyzed,
		fStats.posBroke)

	return bestMove.move
}

func search(game Connect4Game, depth int, moves []int) (s []moveRating) {
	isMaximizingPlayer := game.PlrTurn == CPlr1Max
	moveRatingCh := make(chan moveRating)

	for _, moveCol := range moves {
		go func(m int) {
			posEval := minMax(
				cPutPieceOnBoard(game.Board, m, game.PlrTurn),
				depth,
				-highNumber,
				highNumber,
				// opposite of max because is placing now the other one
				!isMaximizingPlayer)

			moveRatingCh <- moveRating{move: m, eval: posEval}
		}(moveCol)
	}

	numOfPossible := len(moves)
	for i := 0; i < numOfPossible; i++ { // for each goroutine started, get its result
		t := <-moveRatingCh

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
			return highNumber + depth
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
