package game

const highNumber int = 999990000

const (
	CWonS              = highNumber // won/lost
	CCenterS           = 1023       // col 4
	CMidS              = 487        // col 3/5
	CWinnableConnect3S = 734
	CWinnableConnect2S = 304
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
