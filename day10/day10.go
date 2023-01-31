package day10

import (
	"fmt"
	"io"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
)

type State struct {
	cycle    int
	register int
	total    int       // part 1
	out      io.Writer // part 2
}

const (
	// Part 1
	logStart = 20
	logStep  = 40

	// Part 2
	width = 40
	// height = 6

	spriteLen    = 3
	spriteRadius = spriteLen / 2

	pixelOn  = "#"
	pixelOff = "."
)

func render(s *State) {
	s.cycle++

	if (s.cycle-logStart)%logStep == 0 {
		// For part 1
		s.total += s.cycle * s.register
	}

	hPos := (s.cycle - 1) % width

	if hPos >= s.register-spriteRadius && hPos <= s.register+spriteRadius {
		_, err := fmt.Fprint(s.out, pixelOn)
		if err != nil {
			errexit.HandleMainError(err)
		}
	} else {
		_, err := fmt.Fprint(s.out, pixelOff)
		if err != nil {
			errexit.HandleMainError(err)
		}
	}

	if hPos == width-1 {
		_, err := fmt.Fprintln(s.out)
		if err != nil {
			errexit.HandleMainError(err)
		}
	}
}

func Solve(inputPath string, part int) any {
	s := &State{
		cycle:    0,
		register: 1,
		total:    0,
	}
	if part == 1 {
		s.out = io.Discard
	} else { // part == 2
		s.out = &strings.Builder{}
	}

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		switch {
		case strings.HasPrefix(line, "noop"):
			render(s)

		case strings.HasPrefix(line, "addx "):
			count := parse.Int(line[len("addx "):])

			render(s)
			render(s)
			s.register += count
		default:
			errexit.HandleMainError(fmt.Errorf("invalid command %v", line))
		}
	})

	if part == 1 {
		return s.total
	} else { // part == 2
		return s.out.(*strings.Builder).String()
	}
}
