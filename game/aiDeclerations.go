package game

import "time"

const highNumber int = 100_000_000_000

const (
	CWonS              = highNumber // high score for winning
	CWinnableConnect3S = 92         // score for winnable connect-3 sequence
	CWinnableConnect2S = 36         // score for winnable connect-2 sequence
	CHeatS             = 10
)

var locationHeatMap = [CRowsNum][CColsNum]int{
	{0, 1, 3, 8, 3, 1, 0},
	{1, 2, 6, 10, 6, 2, 1},
	{1, 3, 8, 12, 8, 3, 1},
	{1, 3, 8, 12, 8, 3, 1},
	{1, 2, 6, 10, 6, 2, 1},
	{0, 1, 3, 8, 3, 1, 0}}

type flameStats struct {
	posAnalyzed int
	posBroke    int
}

var fStats = &flameStats{}

type moveRating struct {
	move int
	eval int
}

const timeoutTime = 800 * time.Millisecond
