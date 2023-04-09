package game

func staticEval(board *CBoard) (score int) {
	score += getConnectionScore(board)
	score += getLocationScore(board)

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

/* func getHeatMapScore(board *CBoard) (plr1V, plr2V int) {
	// score += getConnectionScore(board)
	// plr1V, plr2V := getHeatMapScore(board)
	// score += (plr1V - plr2V) * CHeatS
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
} */
