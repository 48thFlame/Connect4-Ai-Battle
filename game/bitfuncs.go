package game

import "fmt"

type BITMASK uint64

const BITS = CRowsNum * CColsNum

var PowersArray [BITS + 1]BITMASK // array of size 43 to hold the powers of 2 from 0 to 42
var MasksMap map[int][]BITMASK
var AllCombinations []BITMASK
var inited = false

var BOTTOMINDEX = [7]BITMASK{35, 36, 37, 38, 39, 40, 41}

type BitBoard struct {
	player1 BITMASK
	player2 BITMASK
}

func (board *BitBoard) getFullMask() BITMASK {
	return board.player1 | board.player2
}

func (board *BitBoard) GetPlayerMask(player CPlr) BITMASK {
	if player == CPlr1Max {
		return board.player1
	}
	return board.player2
}

func (board *BitBoard) getSpotVal(row, col int) CPlr {
	if board.player1&PowersArray[getInd(row, col)] > 0 {
		return CPlr1Max
	}
	if board.player2&PowersArray[getInd(row, col)] > 0 {
		return CPlr2Min
	}
	return CNone

}

func (board *BitBoard) setPlayerMask(player CPlr, mask BITMASK) {
	if player == CPlr1Max {
		board.player1 |= mask
	} else {
		board.player2 |= mask
	}
}

func (board *BitBoard) GetAvailableMoves() []int {
	moves := []int{}
	b := board.player1 | board.player2
	for i := 0; i < CColsNum; i++ {
		if b&PowersArray[i] == 0 {
			moves = append(moves, i)
		}
	}

	return moves
}

func (board *BitBoard) PlayMove(move int, player CPlr) int {
	curMask := board.getFullMask()
	for row := CRowsNum - 1; row >= 0; row-- {
		power := PowersArray[getInd(row, move)]
		// ind := BOTTOMINDEX[move] - BITMASK(row*CColsNum)
		if curMask&power == 0 {
			board.setPlayerMask(player, power)
			return row
		}
	}
	return -1
}

func (board *BitBoard) CheckIfWinAtSpot(player CPlr, row, col int) bool {
	return bitCheckIfWin(board, player, MasksMap[getInd(row, col)])
}

func (board *BitBoard) CheckIfWinFullBoard(player CPlr) bool {
	return bitCheckIfWin(board, player, AllCombinations)
}

func BitDumpMask(mask BITMASK) {
	fmt.Println("------MASK-------")
	for row := 0; row < CRowsNum; row++ {
		for col := 0; col < CColsNum; col++ {
			val := 0
			if mask&PowersArray[getInd(row, col)] > 0 {
				val = 1
			}
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
	fmt.Println("------------------")

}

func BitDumpCombinations(masks []BITMASK) {
	for i, m := range masks {
		fmt.Printf("Mask %d:\n", i)
		BitDumpMask(m)
	}

}

func (board *BitBoard) Dump() {
	fmt.Println("----Board---------")
	for row := 0; row < CRowsNum; row++ {
		for col := 0; col < CColsNum; col++ {
			val := board.getSpotVal(row, col)
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
	fmt.Println("------------------")

}

func InitBits() {
	if !inited {
		Powers()
		genCheckPatterns()

		inited = true
		fmt.Println("Bits initialized")
	}

}

func Powers() {
	for i := 0; i <= BITS; i++ {
		PowersArray[i] = 1 << BITMASK(i) // shift 1 to the left i times to calculate the power of 2
	}
	// fmt.Println(PowersArray) // print the array of powers of 2

}

func genMaskCombinations(rowNum, colNum int) []BITMASK {
	/*
		[[2 0 0 0 2 2 2]
		 [1 1 1 0 1 0 1]
		 [1 2 2 0 0 0 2]
		 [2 1 1 1 1 1 1]
		 [1 0 0 2 2 2 2]
		 [0 0 1 1 1 1 2]]
	*/
	var combinations []BITMASK

	// row combinations
	for rowI := 0; rowI < rowNum; rowI++ {
		// for rowI, _ := range board {
		for colI := 0; colI < colNum-3; colI++ {
			combination := PowersArray[getInd(rowI, colI)] | PowersArray[getInd(rowI, colI+1)] | PowersArray[getInd(rowI, colI+2)] | PowersArray[getInd(rowI, colI+3)]
			combinations = append(combinations, combination)
		}
	}

	// col combinations
	for col := 0; col < colNum; col++ {
		for rowI := 0; rowI < rowNum-3; rowI++ {
			combination := PowersArray[getInd(rowI, col)] | PowersArray[getInd(rowI+1, col)] | PowersArray[getInd(rowI+2, col)] | PowersArray[getInd(rowI+3, col)]
			combinations = append(combinations, combination)
		}
	}

	for rowI := 0; rowI < rowNum-3; rowI++ {
		for colI := 0; colI < colNum-3; colI++ {
			combination := PowersArray[getInd(rowI, colI)] | PowersArray[getInd(rowI+1, colI+1)] | PowersArray[getInd(rowI+2, colI+2)] | PowersArray[getInd(rowI+3, colI+3)]
			// combination := []CPlr{board[rowI][colI], board[rowI+1][colI+1], board[rowI+2][colI+2], board[rowI+3][colI+3]}
			combinations = append(combinations, combination)
		}
	}

	// iterate over every diagonal (starting from top-right corner)
	for rowI := 0; rowI < rowNum-3; rowI++ {
		for j := 3; j < 7; j++ {
			combination := PowersArray[getInd(rowI, j)] | PowersArray[getInd(rowI+1, j-1)] | PowersArray[getInd(rowI+2, j-2)] | PowersArray[getInd(rowI+3, j-3)]
			// combination := []CPlr{board[rowI][j], board[rowI+1][j-1], board[rowI+2][j-2], board[rowI+3][j-3]}
			combinations = append(combinations, combination)
		}
	}

	return combinations
}

func genCheckPatterns() {
	AllCombinations = genMaskCombinations(CRowsNum, CColsNum)

	MasksMap = make(map[int][]BITMASK, BITS)

	for row := 0; row < CRowsNum; row++ {
		for col := 0; col < CColsNum; col++ {
			MasksMap[getInd(row, col)] = getMasksForSpot(AllCombinations, row, col)
		}
	}

}

func getInd(row, col int) int {
	return row*CColsNum + col
}

func getMasksForSpot(combinations []BITMASK, row, col int) []BITMASK {
	power := PowersArray[getInd(row, col)]
	var masks []BITMASK

	for _, combo := range combinations {
		if (power & combo) > 0 {
			masks = append(masks, combo)
		}
	}

	return masks
}

func convertBoardToBits(board *CBoard) BitBoard {
	bboard := BitBoard{
		player1: 0,
		player2: 0,
	}

	for row := 0; row < CRowsNum; row++ {
		for col := 0; col < CColsNum; col++ {
			switch board[row][col] {
			case CPlr1Max:
				bboard.player1 |= PowersArray[getInd(row, col)]
			case CPlr2Min:
				bboard.player2 |= PowersArray[getInd(row, col)]
			}
		}
	}

	return bboard
}

func bitCheckIfWin(board *BitBoard, player CPlr, combinations []BITMASK) bool {
	// fmt.Printf("Checking %d combinations\n", len(combinations))
	// BitDumpCombinations(combinations)

	for _, mask := range combinations {
		// BitDumpMask(mask)
		// res := mask & board.GetPlayerMask(player)
		// BitDumpMask(res)
		if mask == mask&board.GetPlayerMask(player) {
			return true
		}
	}

	return false
}

// func bitCheckIfWinFullBoard(board *CBoard, player CPlr) bool {
// 	return bitCheckIfWin(board, player, AllCombinations)
// }
