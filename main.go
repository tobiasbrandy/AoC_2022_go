package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/day1"
	"github.com/tobiasbrandy/AoC_2022_go/day10"
	"github.com/tobiasbrandy/AoC_2022_go/day11"
	"github.com/tobiasbrandy/AoC_2022_go/day12"
	"github.com/tobiasbrandy/AoC_2022_go/day13"
	"github.com/tobiasbrandy/AoC_2022_go/day14"
	"github.com/tobiasbrandy/AoC_2022_go/day15"
	"github.com/tobiasbrandy/AoC_2022_go/day17"
	"github.com/tobiasbrandy/AoC_2022_go/day2"
	"github.com/tobiasbrandy/AoC_2022_go/day3"
	"github.com/tobiasbrandy/AoC_2022_go/day4"
	"github.com/tobiasbrandy/AoC_2022_go/day5"
	"github.com/tobiasbrandy/AoC_2022_go/day6"
	"github.com/tobiasbrandy/AoC_2022_go/day7"
	"github.com/tobiasbrandy/AoC_2022_go/day8"
	"github.com/tobiasbrandy/AoC_2022_go/day9"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"time"
)

type AoCSolver func(string, int) any

func NotImplementedSolver() AoCSolver {
	errexit.HandleMainError(errors.New("not implemented"))
	return nil
}

func PartsSolver(part1, part2 func(string) any) AoCSolver {
	return func(inputPath string, part int) any {
		switch part {
		case 1:
			return part1(inputPath)
		case 2:
			return part2(inputPath)
		default:
			panic("unreachable")
		}
	}
}

var DaySolvers = []AoCSolver{
	day1.Solve,
	day2.Solve,
	PartsSolver(day3.Part1, day3.Part2),
	day4.Solve,
	day5.Solve,
	day6.Solve,
	day7.Solve,
	PartsSolver(day8.Part1, day8.Part2),
	day9.Solve,
	day10.Solve,
	day11.Solve,
	day12.Solve,
	PartsSolver(day13.Part1, day13.Part2),
	day14.Solve,
	PartsSolver(day15.Part1, day15.Part2),
	NotImplementedSolver(), // Day 16
	day17.Solve,
}

func main() {
	day := flag.Int("day", 0, "AoC challenge day number.")
	part := flag.Int("part", 1, "AoC challenge part number. Default: 1.")
	inputPath := flag.String("input", "", "Path to the input file. Default: `day{day}/input.txt`.")
	takeTime := flag.Bool("time", false, "Print execution time")
	test := flag.Bool("test", false, "Ignore `input` parameter and use `day{day}/test.txt` as input.")
	flag.Parse()

	if *day < 1 || *day > len(DaySolvers) {
		errexit.HandleArgsError(fmt.Errorf("day must be between 1 and %d: not %d", len(DaySolvers), *day))
	}

	if *part != 1 && *part != 2 {
		errexit.HandleArgsError(fmt.Errorf("AoC challenges only have part 1 or 2, not part %d", *part))
	}

	if *test {
		*inputPath = fmt.Sprintf("day%d/test.txt", *day)
	} else if *inputPath == "" {
		*inputPath = fmt.Sprintf("day%d/input.txt", *day)
	}

	t := time.Now()
	ret := DaySolvers[*day-1](*inputPath, *part)
	execTime := time.Since(t)

	fmt.Println(ret)
	if *takeTime {
		fmt.Println("Execution time:", execTime)
	}
}
