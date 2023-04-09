package game

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const ITERATIONS = 50000

func YtAiGetMove(game Connect4Game) (col int) {

	InitBits()
	// fmt.Println("Dumping")
	// dumpBoard(&game.Board)
	fmt.Printf("Generating move for %d\n", game.PlrTurn)

	bitBoard := convertBoardToBits(&game.Board)

	moves := bitBoard.GetAvailableMoves()

	winning, wmove := getWinningMove(&bitBoard, game.PlrTurn, &moves)
	if winning {
		fmt.Println("Direct Win")
		return wmove
	}

	winning, wmove = getWinningMove(&bitBoard, cGetOtherPlayer(game.PlrTurn), &moves)
	if winning {
		fmt.Println("Forced Block")
		return wmove
	}

	//col = moves[rand.Intn(len(moves))]
	col = MonteCarloAI(&bitBoard, game.PlrTurn, ITERATIONS)

	return col
}

func getWinningMove(board *BitBoard, player CPlr, moves *[]int) (bool, int) {
	winning := getFirstFilteredMove(board, player, moves, isWinningMove)
	if winning > 0 {
		return true, winning
	}
	return false, -1
}

func isMovePossible(board *CBoard, col int) (possible bool) {
	for i := 0; i < CRowsNum; i++ { // start at top and f sees an empty spot there is at least one possible room
		if board[i][col] == CNone {
			return true
		}
	}

	return false
}

func getAvailableMoves(board *CBoard) []int {
	a := []int{}

	for i := 0; i < CColsNum; i++ {
		if isMovePossible(board, i) {
			a = append(a, i)
		}
	}

	return a
}

func isWinningMove(board BitBoard, move int, player CPlr) bool {
	// b := game.place(move+1, player)
	row := board.PlayMove(move, player)

	// dumpBoard(&board)

	return board.CheckIfWinAtSpot(player, row, move)
}

func getWinScore(winplayer CPlr, maxplayer CPlr) int {
	if maxplayer == winplayer {
		// fmt.Printf("%d won so returning: 1\n", winplayer)
		return 1
	}
	// fmt.Printf("%d won so returning: -1\n", winplayer)

	return -1
}

func runRandomGame(board BitBoard, maxplayer CPlr, nextmove int) int {
	row := board.PlayMove(nextmove, maxplayer)
	gameover := board.CheckIfWinAtSpot(maxplayer, row, nextmove)

	nextplayer := cGetOtherPlayer(maxplayer)

	for !gameover {
		// fmt.Println("Begining Random Section")
		moves := board.GetAvailableMoves()
		if len(moves) == 0 {
			return 0
		}

		winning, _ := getWinningMove(&board, nextplayer, &moves)
		if winning {
			return getWinScore(nextplayer, maxplayer)
		}

		winning, move := getWinningMove(&board, cGetOtherPlayer(nextplayer), &moves)
		if !winning {
			move = moves[rand.Intn(len(moves))]
		}
		// move := moves[rand.Intn(len(moves))]
		board.PlayMove(move, nextplayer)
		// if board.CheckIfWinAtSpot(nextplayer, row, move) {
		if board.CheckIfWinFullBoard(nextplayer) {
			return getWinScore(nextplayer, maxplayer)
		}

		nextplayer = cGetOtherPlayer(nextplayer)
	}
	return 1
}

func playPiece(board *CBoard, col int, piece CPlr) int {
	for i := CRowsNum - 1; i >= 0; i-- { // going from 6-0 to start at bottom
		if board[i][col] == CNone {
			board[i][col] = piece
			return i
		}
	}

	return 7

}

func myMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func validCoords(row int, col int) bool {
	return (row >= 0) && (row < CRowsNum) && (col >= 0) && (col < CColsNum)
}

func rowincer(row int, col int) (int, int) {
	return row, col + 1
}

func colincer(row int, col int) (int, int) {
	return row - 1, col
}

func diagposincer(row int, col int) (int, int) {
	return row - 1, col + 1
}

func diagnegincer(row int, col int) (int, int) {
	return row + 1, col + 1
}

func checkConnect(n int, board *CBoard, lookingFor CPlr, row int, col int, incer func(int, int) (int, int)) bool {
	counter := 0

	for validCoords(row, col) {
		// fmt.Printf("%d %d\n", row, col)
		if board[row][col] == lookingFor {
			counter++
			if counter >= n {
				return true
			}
		} else {
			counter = 0
		}

		row, col = incer(row, col)
	}

	return false
}

func checkForWinSlow(board *CBoard, player CPlr, row int, col int) (won bool) {
	// row checks
	// fmt.Printf("Move was %d %d \n", row, col)
	// fmt.Println("Check Row")
	won = checkConnect(4, board, player, row, 0, rowincer)
	if won {
		return won
	}
	// fmt.Println("Check Col")
	won = checkConnect(4, board, player, 5, col, colincer)
	if won {
		return won
	}

	delta := myMin((CRowsNum-1)-row, col)
	posrow := row + delta
	poscol := col - delta
	// fmt.Println("Check Diag Pos")

	won = checkConnect(4, board, player, posrow, poscol, diagposincer)
	if won {
		return won
	}

	delta = myMin(row, col)
	negrow := row - delta
	negcol := col - delta
	// fmt.Println("Check Diag Neg")

	won = checkConnect(4, board, player, negrow, negcol, diagnegincer)
	if won {
		return won
	}

	return false
}

func getMovesWithCondition(board *CBoard, player CPlr, filter func(*CBoard, int, CPlr) bool) []int {
	a := []int{}

	for i := 1; i < CColsNum+1; i++ {
		if cMovePossible(*board, i) && filter(board, i, player) {
			a = append(a, i)
		}
	}

	return a
}

func getFirstFilteredMove(board *BitBoard, player CPlr, moves *[]int, filter func(BitBoard, int, CPlr) bool) int {

	for i := range *moves {
		if filter(*board, i, player) {
			return i
		}
	}

	return -1
}

func filterMoves(board *CBoard, player CPlr, moves *[]int, filter func(CBoard, int, CPlr) bool) []int {
	a := []int{}

	for i := range *moves {
		if filter(*board, i, player) {
			a = append(a, i)
		}
	}

	return a
}

func MonteCarloAI(board *BitBoard, aiplayer CPlr, iterations int) int {
	// perform Monte Carlo simulation for the given number of iterations
	// and return the best move to play

	moves := board.GetAvailableMoves()
	numMoves := len(moves)
	// fmt.Println(moves)

	if numMoves == 0 {
		return 7
	}

	var wg sync.WaitGroup
	wg.Add(numMoves)

	scores := make([]int, numMoves)
	for i, move := range moves {
		go func(index int, m int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// fmt.Printf("Running Monte Carlo Iteration %d\n", j)

				scores[index] += runRandomGame(*board, aiplayer, m)
				// evaluate the final game state and update the score for this move
			}
		}(i, move)

	}

	wg.Wait()

	// choose the move with the highest average score
	bestScore := math.MinInt
	bestMove := moves[0]
	fmt.Println(scores)
	for i, score := range scores {
		if score > bestScore {
			bestScore = score
			bestMove = moves[i]
		}
	}

	return bestMove
}

func dumpBoard(board *CBoard) {
	fmt.Println("Board:")
	for r := range board {
		fmt.Println(board[r])
	}
}

func evaluateBoard(board *CBoard) int64 {

	return 0
}

// func MaxMin(currentState GameState, depth int, maximizingPlayer bool) float64 {
// 	// Check if the game is over or if the depth limit has been reached
// 	if gameOver(currentState) || depth == 0 {
// 		return evaluateState(currentState)
// 	}

// 	if maximizingPlayer {
// 		bestValue := math.Inf(-1)
// 		for _, nextState := range generateNextStates(currentState) {
// 			value := MaxMin(nextState, depth-1, false)
// 			bestValue = math.Max(bestValue, value)
// 		}
// 		return bestValue
// 	} else {
// 		bestValue := math.Inf(1)
// 		for _, nextState := range generateNextStates(currentState) {
// 			value := MaxMin(nextState, depth-1, true)
// 			bestValue = math.Min(bestValue, value)
// 		}
// 		return bestValue
// 	}
// }

// func findBestMove(currentState GameState, depth int) GameState {
// 	bestValue := math.Inf(-1)
// 	var bestMove GameState

// 	for _, nextState := range generateNextStates(currentState) {
// 		value := MaxMin(nextState, depth-1, false)
// 		if value > bestValue {
// 			bestValue = value
// 			bestMove = nextState
// 		}
// 	}

// 	return bestMove
// }
