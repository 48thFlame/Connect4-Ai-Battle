package game

import (
	"math/rand"
)

func RandAiGetMove(game Connect4Game) (col int) {
	a := cGetAvailableMoves(game.Board)
	col = a[rand.Intn(len(a))]

	return col
}
