package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/48thFlame/Connect4-Ai-Battle/game"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type matchResult int

const (
	flame1Draws matchResult = iota
	yt1Draws

	flame1Wins
	flame2Wins
	yt1Wins
	yt2Wins
)

type matchScores struct {
	flame1Draws int
	yt1Draws    int
	flame1Wins  int
	flame2Wins  int
	yt1Wins     int
	yt2Wins     int
}

func main() {
	connect4 := game.NewConnect4Game()

	// gameLoop(connect4, game.FlameAiGetMove, game.FlameAiGetMove)
	// gameLoop(connect4, game.RandAiGetMove, game.FlameAiGetMove)
	// gameLoop(connect4, true, game.HumanGetMove, game.FlameAiGetMove)
	// gameLoop(connect4, true, game.YtAiGetMove, game.FlameAiGetMove)
	gameLoop(connect4, game.FlameAiGetMove, game.YtAiGetMove)
	return
	numOfGames := 12
	bigR := &matchScores{}

	for i := 0; i < numOfGames; i++ {
		fmt.Printf("Game %v starting\n", i)
		flameFirst := i%2 == 0
		fmt.Printf("Is flame first: %v\n", flameFirst)
		result := simulateMatch(flameFirst)
		updateBigR(bigR, result)
		fmt.Println("Result:", result)
		fmt.Println()
		fmt.Printf("%#v", bigR)
		fmt.Println()
		fmt.Println("#########")
		fmt.Println()
		fmt.Println()
	}

	fmt.Println("Results are in....")
	fmt.Printf("%#v", bigR)
	fmt.Println()

}

func updateBigR(bigR *matchScores, r matchResult) {
	switch r {
	case flame1Draws:
		bigR.flame1Draws++

	case yt1Draws:
		bigR.yt1Draws++

	case flame1Wins:
		bigR.flame1Wins++

	case flame2Wins:
		bigR.flame2Wins++

	case yt1Wins:
		bigR.yt1Wins++

	case yt2Wins:
		bigR.yt2Wins++

	}
}

func simulateMatch(flameFirst bool) matchResult {
	c := game.NewConnect4Game()
	var plr1Input, plr2input func(game game.Connect4Game) int

	if flameFirst {
		plr1Input = game.FlameAiGetMove
		plr2input = game.YtAiGetMove
	} else {
		plr1Input = game.YtAiGetMove
		plr2input = game.FlameAiGetMove
	}

	winner := gameLoop(c, plr1Input, plr2input)

	if flameFirst {
		switch winner {
		case game.CNone:
			return flame1Draws
		case game.CPlr1Max:
			return flame1Wins
		case game.CPlr2Min:
			return yt2Wins
		}
	} else {
		switch winner {
		case game.CNone:
			return yt1Draws
		case game.CPlr1Max:
			return yt1Wins
		case game.CPlr2Min:
			return flame2Wins
		}
	}

	return flame1Draws
}

func gameLoop(connect4 *game.Connect4Game, plr1Input, plr2Input func(game game.Connect4Game) int) game.CPlr {
loop:
	for {
		fmt.Println("-----------------------")

		isPlr1Turn := connect4.PlrTurn == game.CPlr1Max

		var plrCol int
		var plrString string

		if isPlr1Turn {
			plrCol = plr1Input(*connect4)
			plrString = "Player1"
		} else {
			plrCol = plr2Input(*connect4)
			plrString = "Player2"
		}

		fmt.Printf("Game turnNumber: %v\n", connect4.TurnNum)
		fmt.Printf("%v went at col: %v\n", plrString, plrCol+1)

		good := connect4.Turn(plrCol)
		if !good {
			fmt.Printf("game: %#v\n", connect4)
			break
		}

		fmt.Println(connect4GameToString(connect4))

		switch connect4.GameState {
		case game.CStatePlaying:
			continue
		case game.CStateDraw:
			fmt.Println("Game ended in a draw.")

			break loop

		case game.CStatePlr1Won:
			fmt.Println("Player1 won!")

			return game.CPlr1Max

		case game.CStatePlr2Won:
			fmt.Println("Player2 won!")

			return game.CPlr2Min
		}
	}

	return game.CNone
}

func connect4GameToString(game *game.Connect4Game) string {
	var sb strings.Builder

	for _, row := range game.Board {
		for _, spot := range row {
			sb.WriteString(cPlrToString(spot))
		}
		sb.WriteRune('\n')
	}

	sb.WriteString("1️⃣ 2️⃣ 3️⃣ 4️⃣ 5️⃣ 6️⃣ 7️⃣")

	return sb.String()
}

func cPlrToString(c game.CPlr) string {
	switch c {
	case game.CNone:
		return "🔳"
	case game.CPlr1Max:
		return "🔵"
	case game.CPlr2Min:
		return "🔴"
	default:
		return ""
	}
}
