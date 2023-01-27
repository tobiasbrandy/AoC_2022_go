package main

import (
	"errors"
	"flag"
	"fmt"
	day1 "github.com/tobiasbrandy/AoC_2022_go/Day1"
	day10 "github.com/tobiasbrandy/AoC_2022_go/Day10"
	day11 "github.com/tobiasbrandy/AoC_2022_go/Day11"
	day12 "github.com/tobiasbrandy/AoC_2022_go/Day12"
	day13 "github.com/tobiasbrandy/AoC_2022_go/Day13"
	day14 "github.com/tobiasbrandy/AoC_2022_go/Day14"
	day15 "github.com/tobiasbrandy/AoC_2022_go/Day15"
	day17 "github.com/tobiasbrandy/AoC_2022_go/Day17"
	day2 "github.com/tobiasbrandy/AoC_2022_go/Day2"
	day3 "github.com/tobiasbrandy/AoC_2022_go/Day3"
	day4 "github.com/tobiasbrandy/AoC_2022_go/Day4"
	day5 "github.com/tobiasbrandy/AoC_2022_go/Day5"
	day6 "github.com/tobiasbrandy/AoC_2022_go/Day6"
	day7 "github.com/tobiasbrandy/AoC_2022_go/Day7"
	day8 "github.com/tobiasbrandy/AoC_2022_go/Day8"
	day9 "github.com/tobiasbrandy/AoC_2022_go/Day9"
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
	inputPath := flag.String("input", "", "Path to the input file. Default: `Day{day}/input.txt`.")
	takeTime := flag.Bool("time", false, "Print execution time")
	test := flag.Bool("test", false, "Ignore `input` parameter and use `Day{day}/test.txt` as input.")

	flag.Parse()

	if *day < 1 || *day > len(DaySolvers) {
		errexit.HandleArgsError(fmt.Errorf("day must be between 1 and %d: not %d", len(DaySolvers), *day))
	}

	if *part != 1 && *part != 2 {
		errexit.HandleArgsError(fmt.Errorf("AoC challenges only have part 1 or 2, not part %d", *part))
	}

	if *test {
		*inputPath = fmt.Sprintf("Day%d/test.txt", *day)
	} else if *inputPath == "" {
		*inputPath = fmt.Sprintf("Day%d/input.txt", *day)
	}

	t := time.Now()
	ret := DaySolvers[*day-1](*inputPath, *part)
	execTime := time.Since(t)

	fmt.Println(ret)
	if *takeTime {
		fmt.Println("Execution time:", execTime)
	}
}
