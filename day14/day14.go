package day14

import (
	"fmt"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

func parseNode(node string) pos.D2 {
	coords := strings.Split(node, ",")
	return pos.New2D(
		parse.Int(coords[0]),
		parse.Int(coords[1]),
	)
}

func Solve(inputPath string, part int) any {
	source := pos.New2D(500, 0)
	rockSet := set.Set[pos.D2]{}
	maxY := 0

	// Build initial walls
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		nodes := strings.Split(line, " -> ")
		prev := parseNode(nodes[0])
		rockSet.Add(prev)
		if prev.Y > maxY {
			maxY = prev.Y
		}

		for _, node := range nodes[1:] {
			curr := parseNode(node)
			if curr.Y > maxY {
				maxY = curr.Y
			}

			if curr.X == prev.X {
				dir := mathext.Sign(curr.Y - prev.Y)
				for i := prev.Y + dir; i != curr.Y; i += dir {
					rockSet.Add(pos.New2D(curr.X, i))
				}
			} else if curr.Y == prev.Y {
				dir := mathext.Sign(curr.X - prev.X)
				for i := prev.X + dir; i != curr.X; i += dir {
					rockSet.Add(pos.New2D(i, curr.Y))
				}
			} else {
				errexit.HandleMainError(fmt.Errorf("rock path not in a straight line: prev=%v, curr=%v", prev, curr))
			}

			rockSet.Add(curr)
			prev = curr
		}
	})

	if part == 2 {
		maxY += 2 // Set floor y position
	}

	// Simulate
	totalSand := 0
	for {
		curr := source

		for curr.Y < maxY {
			if part == 2 && curr.Y+1 == maxY {
				// Sand hit floor
				rockSet.Add(curr)
				break
			} else if test := pos.New2D(curr.X, curr.Y+1); !rockSet.Contains(test) {
				curr = test
			} else if test := pos.New2D(curr.X-1, curr.Y+1); !rockSet.Contains(test) {
				curr = test
			} else if test := pos.New2D(curr.X+1, curr.Y+1); !rockSet.Contains(test) {
				curr = test
			} else {
				// Put sand to rest as rock
				rockSet.Add(curr)
				break
			}
		}

		infiniteFlow := curr.Y >= maxY
		sourceBlocked := curr == source
		if infiniteFlow || sourceBlocked {
			if sourceBlocked {
				totalSand++
			}
			break
		}

		totalSand++
	}

	return totalSand
}
