package day2

import (
	"errors"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
)

type RPSMove int

const (
	Rock RPSMove = iota + 1
	Paper
	Scissor
)

var rpsLoses = map[RPSMove]RPSMove{
	Rock:    Paper,
	Paper:   Scissor,
	Scissor: Rock,
}

var rpsWins = map[RPSMove]RPSMove{
	Rock:    Scissor,
	Paper:   Rock,
	Scissor: Paper,
}

func (move RPSMove) Beats(other RPSMove) bool {
	return rpsWins[move] == other
}

func (move RPSMove) WinsTo() RPSMove {
	return rpsWins[move]
}

func (move RPSMove) LosesTo() RPSMove {
	return rpsLoses[move]
}

func Solve(inputPath string, part int) any {
	if part == 1 {
		errexit.HandleMainError(errors.New("part 1 not implemented :("))
	}

	score := 0

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		other := RPSMove((line[0] - 'A') + 1)

		var you RPSMove
		switch line[2] {
		case 'X':
			you = other.WinsTo()
		case 'Y':
			you = other
		case 'Z':
			you = other.LosesTo()
		}

		// Score by choice
		score += int(you)

		// Score by result
		if you.Beats(other) { // You won
			score += 6
		} else if you == other { // Tie
			score += 3
		}
		// Else you lose => score += 0
	})

	return score
}
