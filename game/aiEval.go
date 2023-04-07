package game

func staticEval(board *CBoard) (score int) {
	score += getConnectionScore(board)
	plr1V, plr2V := getHeatMapScore(board)

	score += (plr1V - plr2V) * CHeatS

	return score
}

func getHeatMapScore(board *CBoard) (plr1V, plr2V int) {
	for rowI, row := range board {
		for colI, spot := range row {
			if spot == CPlr1Max {
				plr1V += locationHeatMap[rowI][colI]
			} else if spot == CPlr2Min {
				plr2V += locationHeatMap[rowI][colI]
			}
		}
	}

	return plr1V, plr2V
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
