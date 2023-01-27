package day6

import (
	"fmt"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

func isUnique[T comparable](slice []T) bool {
	l := len(slice)
	s := set.New[T](l)
	s.AddAll(slice)
	return s.Len() == l
}

func Solve(inputPath string, part int) any {
	var uniqueLen int
	if part == 1 {
		uniqueLen = 4
	} else { // part == 2
		uniqueLen = 14
	}

	var ret int
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		if len(line) < uniqueLen {
			errexit.HandleMainError(fmt.Errorf("line is smaller than the required unique characters: %d", uniqueLen))
			return
		}

		for i := range line[4:] {
			if isUnique([]rune(line[i : i+uniqueLen])) {
				ret = i + uniqueLen
				return
			}
		}

		errexit.HandleMainError(fmt.Errorf("no unique set of %d characters", uniqueLen))
	})

	return ret
}
