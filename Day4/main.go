package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"flag"
	"fmt"
	"strconv"
	"strings"
)

func split2(s string, sep rune) (string, string) {
	sepIdx := strings.IndexRune(s, sep)
	return s[:sepIdx], s[sepIdx+1:]
}

func solve(filePath string, part int) {
	total := 0

	internal.ForEachFileLine(filePath, internal.HandleScanError, func(line string) {
		int1, int2 := split2(line, ',')
		int1l, int1r := split2(int1, '-')
		int2l, int2r := split2(int2, '-')

		l1, _ := strconv.Atoi(int1l)
		r1, _ := strconv.Atoi(int1r)
		l2, _ := strconv.Atoi(int2l)
		r2, _ := strconv.Atoi(int2r)

		if part == 1 {
			// Fully overlap
			if (l1 >= l2 && r1 <= r2) || (l2 >= l1 && r2 <= r1) {
				total++
			}
		} else { // part == 2
			// Overlap at all
			if (l1 <= r2 && l1 >= l2) || (l2 <= r1 && l2 >= l1) || (r1 >= l2 && r1 <= r2) || (r2 >= l1 && r2 <= r1) {
				total++
			}
		}
	})

	fmt.Println(total)
}

func main() {
	inputPath := flag.String("input", "input.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	if *part != 1 && *part != 2 {
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}

	solve(*inputPath, *part)
}
