package day4

import (
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
)

func split2(s string, sep rune) (string, string) {
	sepIdx := strings.IndexRune(s, sep)
	return s[:sepIdx], s[sepIdx+1:]
}

func Solve(inputPath string, part int) any {
	total := 0

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		int1, int2 := split2(line, ',')
		int1l, int1r := split2(int1, '-')
		int2l, int2r := split2(int2, '-')

		l1 := parse.Int(int1l)
		r1 := parse.Int(int1r)
		l2 := parse.Int(int2l)
		r2 := parse.Int(int2r)

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

	return total
}
