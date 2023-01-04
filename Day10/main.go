package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	width  = 40
	height = 6

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
		fmt.Fprint(s.out, pixelOn)
	} else {
		fmt.Fprint(s.out, pixelOff)
	}

	if hPos == width-1 {
		fmt.Fprintln(s.out)
	}
}

func solve(filePath string, part int) {
	s := &State{
		cycle:    0,
		register: 1,
		total:    0,
	}
	if part == 1 {
		s.out = io.Discard
	} else { // part == 2
		s.out = os.Stdout
	}

	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
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
		fmt.Println(s.total)
	}
}

func main() {
	inputPath := flag.String("input", "input.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	if *part != 1 && *part != 2 {
		errexit.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}

	solve(*inputPath, *part)
}
