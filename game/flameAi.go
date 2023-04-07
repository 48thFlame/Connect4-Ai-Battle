package game

import (
	"fmt"
	"math"
)

const highNumber int = 999990000

const (
	CWonS              = highNumber // won/lost
	CCenterS           = 1023       // col 4
	CMidS              = 487        // col 3/5
	CWinnableConnect3S = 734
	CWinnableConnect2S = 304
)

const ( // index num starts at 0
	cCenterCol = 3
	cMidCol1   = 2
	cMidCol2   = 4
)

type flameStats struct {
	posAnalyzed int
	posBroke    int
}

var fStats = &flameStats{}

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

func getOrderedMoves(board CBoard) []int {
	s := cGetAvailableMoves(board)
	prob := []int{3, 4, 2, 5, 1, 6, 0}

	c := make([]int, 0)

	for _, good := range prob {
		if inSlice(s, good) {
			c = append(c, good)
		}
	}

	return c
}

func minMax(board CBoard, depth int, alpha int, beta int, maximizingPlayer bool) int {
	if gs := cGetGameState(board); gs != CStatePlaying {
		fStats.posAnalyzed++
		switch gs {
		case CStateDraw:
			return 0
		case CStatePlr1Won:
			return (highNumber + depth)
		case CStatePlr2Won:
			return (-highNumber - depth)
		}
	} else if depth == 0 {
		fStats.posAnalyzed++
		return staticEval(&board)
	}

	possibleMoves := getOrderedMoves(board)
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

func staticEval(board *CBoard) (score int) {
	score += getLocationScore(board)
	score += getConnectionScore(board)

	return score
}

func getLocationScore(board *CBoard) (score int) {
	centerPlr1, centerPlr2 := getCenters(board)
	midPl1, midPlr2 := getMiddles(board)
	score += (centerPlr1 - centerPlr2) * CCenterS
	score += (midPl1 - midPlr2) * CMidS
	return score
}

func getCenters(board *CBoard) (centerPlr1, centerPlr2 int) {
	for _, row := range board {
		spot := row[cCenterCol]
		if spot == CPlr1Max {
			centerPlr1++
		} else if spot == CPlr2Min {
			centerPlr2++
		}
	}

	return centerPlr1, centerPlr2
}

func getMiddles(board *CBoard) (midPlr1, midPlr2 int) {
	for _, row := range board {
		spot1 := row[cMidCol1]
		spot2 := row[cMidCol2]

		if spot1 == CPlr1Max || spot2 == CPlr1Max {
			midPlr1++
		} else if spot1 == CPlr2Min || spot2 == CPlr2Min {
			midPlr2++
		}
	}

	return midPlr1, midPlr2
}

func getConnectionScore(board *CBoard) (score int) {
	b := *board
	possibleMoves := getAvailableMoves(board)
	for _, move := range possibleMoves {
		b = cPutPieceOnBoard(b, move, CPlr(-1))
	}
	s := getCombinations(b)

	number3Plr1, number2Plr1 := getConnectionsForPlr(s, CPlr1Max)
	number3Plr2, number2Plr2 := getConnectionsForPlr(s, CPlr2Min)

	score += (number3Plr1 - number3Plr2) * CWinnableConnect3S
	score += (number2Plr1 - number2Plr2) * CWinnableConnect2S

	return score
}

func getConnectionsForPlr(s [][]CPlr, p CPlr) (number3s, number2s int) {
	var winnable2, winnable3 bool

	for _, combination := range s {
		winnable3 = winnableN(combination, p, 3)
		if winnable3 {
			number3s++
			continue
		}

		winnable2 = winnableN(combination, p, 2)
		if winnable2 {
			number2s++
		}
	}

	return number3s, number2s
}
