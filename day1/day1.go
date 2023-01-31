package day1

import (
	"sort"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
)

func sumInt(arr []int) int {
	sum := 0
	for _, elem := range arr {
		sum += elem
	}
	return sum
}

func Solve(inputPath string, part int) any {
	var top int
	if part == 1 {
		top = 1
	} else { // part == 2
		top = 3
	}

	var itemCount []int
	accum := 0

	fileline.ForEachSet(inputPath, errexit.HandleScanError, func(lines []string) {
		for _, line := range lines {
			accum += parse.Int(line)
		}

		itemCount = append(itemCount, accum)
		accum = 0
	})

	sort.Sort(sort.Reverse(sort.IntSlice(itemCount)))
	topSum := sumInt(itemCount[:top])

	return topSum
}
