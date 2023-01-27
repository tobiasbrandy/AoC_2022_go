package day3

import (
	"fmt"
	"unicode"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

func itemPriority(item rune) int {
	// We assume items are ascii letters
	if unicode.IsUpper(item) {
		return int(item-'A') + 27
	} else {
		return int(item-'a') + 1
	}
}

func Part1(inputPath string) any {
	total := 0

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		ruckDiv := len(line) / 2

		ruck1 := line[:ruckDiv]
		ruck2 := line[ruckDiv:]

		ruck1Set := set.Set[rune]{}
		ruck1Set.AddAll([]rune(ruck1))

		for _, item := range ruck2 {
			if ruck1Set.Contains(item) {
				total += itemPriority(item)
				return
			}
		}
	})

	return total
}

func allSetsContain[T comparable](sets []set.Set[T], elem T) bool {
	for _, s := range sets {
		if !s.Contains(elem) {
			return false
		}
	}
	return true
}

func Part2(inputPath string) any {
	groupCount := 3
	total := 0

	fileline.ForEachSetN(inputPath, groupCount, errexit.HandleScanError, func(group []string) {
		if len(group) != groupCount {
			errexit.HandleMainError(fmt.Errorf("input lines are not divisible by %v. Remainder of %v", groupCount, len(group)))
		}

		groupSets := make([]set.Set[rune], groupCount-1)
		for i := 0; i < groupCount-1; i++ {
			s := set.Set[rune]{}
			s.AddAll([]rune(group[i]))
			groupSets[i] = s
		}

		for _, item := range group[groupCount-1] {
			if allSetsContain(groupSets, item) {
				total += itemPriority(item)
				return
			}
		}
	})

	return total
}
