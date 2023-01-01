package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"flag"
	"fmt"
	"unicode"
)

func itemPriority(item rune) int {
	// We assume items are ascii letters
	if unicode.IsUpper(item) {
		return int(item - 'A') + 27
	} else {
		return int(item - 'a') + 1
	}
}

func part1(filePath string) {
	total := 0

	internal.ForEachFileLine(filePath, internal.HandleScanError, func(line string) {
		ruckDiv := len(line)/2

		ruck1 := line[:ruckDiv]
		ruck2 := line[ruckDiv:]

		ruck1Set := internal.Set[rune]{}
		ruck1Set.AddAll([]rune(ruck1))

		for _, item := range ruck2{
			if ruck1Set.Contains(item) {
				total += itemPriority(item)
				return
			}
		}
	})

	fmt.Println(total)
}

func AllSetsContain[T comparable](sets []internal.Set[T], elem T) bool {
	for _, set := range sets {
		if !set.Contains(elem) {
			return false
		}
	}
	return true
}

func part2(filePath string) {
	groupCount := 3
	total := 0

	internal.ForEachFileLineSetN(filePath, groupCount, internal.HandleScanError, func(group []string) {
		if len(group) != groupCount {
			internal.HandleMainError(fmt.Errorf("input lines are not divisible by %v. Remainder of %v", groupCount, len(group)))
		}

		groupSets := make([]internal.Set[rune], groupCount - 1)
		for i := 0; i < groupCount - 1; i++ {
			set := internal.Set[rune]{}
			set.AddAll([]rune(group[i]))
			groupSets[i] = set
		}

		for _, item := range group[groupCount - 1] {
			if AllSetsContain(groupSets, item) {
				total += itemPriority(item)
				return
			}
		}
	})

	fmt.Println(total)
}

func main() {
	inputPath := flag.String("input", "input1.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	switch *part {
	case 1:
		part1(*inputPath)
	case 2:
		part2(*inputPath)
	default:
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}
}
