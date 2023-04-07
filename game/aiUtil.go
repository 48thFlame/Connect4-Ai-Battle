package game

import "sort"

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func inSlice(s []int, a int) bool {
	for _, obj := range s {
		if obj == a {
			return true
		}
	}

	return false
}

func winnableN(s []CPlr, piece CPlr, n int) bool {
	counter := 0

	for _, spot := range s {
		switch spot {
		case piece:
			counter++
			continue
		case CNone:
			continue
		default: // its opponent
			return false
		}
	}

	return counter == n // only equal not bigger so connect 3 don't count also as connect 2
}

func getCombinations(board CBoard) [][]CPlr {
	/*
		[[2 0 0 0 2 2 2]
		 [1 1 1 0 1 0 1]
		 [1 2 2 0 0 0 2]
		 [2 1 1 1 1 1 1]
		 [1 0 0 2 2 2 2]
		 [0 0 1 1 1 1 2]]
	*/
	var combinations [][]CPlr

	// row combinations
	for _, row := range board {
		for colI := 0; colI < 4; colI++ {
			combination := row[colI : colI+4]
			combinations = append(combinations, combination)
		}
	}

	// col combination no because they can just be blocked
	// // col combinations
	// for col := 0; col < CColsNum; col++ {
	// 	for rowI := 0; rowI < 3; rowI++ {
	// 		combination := []CPlr{board[rowI][col], board[rowI+1][col], board[rowI+2][col], board[rowI+3][col]}
	// 		combinations = append(combinations, combination)
	// 	}
	// }

	for rowI := 0; rowI < 3; rowI++ {
		for colI := 0; colI < 4; colI++ {
			combination := []CPlr{board[rowI][colI], board[rowI+1][colI+1], board[rowI+2][colI+2], board[rowI+3][colI+3]}
			combinations = append(combinations, combination)
		}
	}

	// iterate over every diagonal (starting from top-right corner)
	for rowI := 0; rowI < 3; rowI++ {
		for j := 3; j < 7; j++ {
			combination := []CPlr{board[rowI][j], board[rowI+1][j-1], board[rowI+2][j-2], board[rowI+3][j-3]}
			combinations = append(combinations, combination)
		}
	}

	return combinations
}

func getOrderedAvailableMoves(board CBoard) []int {
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

func sortMovesHighLow(s []moveRating) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].eval > s[j].eval
	})
}

func sortMovesLowHigh(s []moveRating) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].eval < s[j].eval
	})
}

func moveRatingsToMoves(s []moveRating) []int {
	r := make([]int, 0)
	for _, m := range s {
		r = append(r, m.move)
	}

	return r
}
