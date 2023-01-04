package main

import (
	"flag"
	"fmt"

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

func part1(inputPath string) {
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

	fmt.Println(score)
}

func main() {
	inputPath := flag.String("input", "input.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	switch *part {
	case 1:
		part1(*inputPath)
	default:
		errexit.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}
}
