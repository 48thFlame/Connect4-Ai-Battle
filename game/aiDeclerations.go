package game

import "time"

const highNumber int = 100_000_000_000

// const (
// 	CWonS              = highNumber // won/lost
// 	CCenterS           = 1023       // col 4
// 	CMidS              = 487        // col 3/5
// 	CWinnableConnect3S = 734
// 	CWinnableConnect2S = 304
// )

const (
	CWonS              = highNumber // high score for winning
	CCenterS           = 1000       // score for center column (col 4)
	CMidS              = 600        // score for middle columns (col 3/5)
	CWinnableConnect3S = 800        // score for winnable connect-3 sequence
	CWinnableConnect2S = 250        // score for winnable connect-2 sequence
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

type moveRating struct {
	move int
	eval int
}

const timeoutTime = 800 * time.Millisecond
