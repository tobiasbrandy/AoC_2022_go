package main

import (
	"flag"
	"fmt"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

func isUnique[T comparable](slice []T) bool {
	len := len(slice)
	set := set.New[T](len)
	set.AddAll(slice)
	return set.Len() == len
}

func solve(filePath string, part int) {
	var uniqueLen int
	if part == 1 {
		uniqueLen = 4
	} else { // part == 2
		uniqueLen = 14
	}

	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
		if len(line) < uniqueLen {
			fmt.Println("Line is smaller than the required unique characters:", uniqueLen)
			return
		}

		for i := range line[4:] {
			if isUnique([]rune(line[i : i+uniqueLen])) {
				fmt.Println(i + uniqueLen)
				return
			}
		}

		fmt.Println("No unique set of", uniqueLen, "characters")
	})
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
