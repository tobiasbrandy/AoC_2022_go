package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"flag"
	"fmt"
	"sort"
	"strconv"
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

	internal.ForEachFileLine(inputPath, internal.HandleScanError, func(line string) {
		if line == "" {
			itemCount = append(itemCount, accum)
			accum = 0
		} else {
			count, err := strconv.Atoi(line)
			if err != nil {
				internal.HandleMainError(err)
			}

			accum += count
		}
	})
	itemCount = append(itemCount, accum)

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
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}
}
