package main

import (
	"flag"
	"fmt"
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

func part1(inputPath string) {
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
	topSum := sumInt(itemCount[:3])

	fmt.Println(topSum)
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
