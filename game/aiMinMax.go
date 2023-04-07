package game

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
