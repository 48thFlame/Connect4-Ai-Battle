package game

import "fmt"

func HumanGetMove(connect4 Connect4Game) int {
	var col int

	fmt.Print("Enter column: ")
	fmt.Scan(&col)

	col--

	possibleMoves := cGetAvailableMoves(connect4.Board)

	if inSlice(possibleMoves, col) {
		return col
	} else {
		fmt.Println("Invalid column!")
		return HumanGetMove(connect4)
	}
}
