package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

type Pos2D struct {
	x, y int
}

func parseNode(node string) Pos2D {
	coords := strings.Split(node, ",")
	return Pos2D{
		parse.Int(coords[0]),
		parse.Int(coords[1]),
	}
}

func solve(filePath string, part int) {
	source := Pos2D{500, 0}
	rockSet := set.Set[Pos2D]{}
	maxY := 0

	// Build initial walls
	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
		nodes := strings.Split(line, " -> ")
		prev := parseNode(nodes[0])
		rockSet.Add(prev)
		if prev.y > maxY {
			maxY = prev.y
		}

		for _, node := range nodes[1:] {
			curr := parseNode(node)
			if curr.y > maxY {
				maxY = curr.y
			}

			if curr.x == prev.x {
				dir := mathext.Sign(curr.y - prev.y)
				for i := prev.y + dir; i != curr.y; i += dir {
					rockSet.Add(Pos2D{curr.x, i})
				}
			} else if curr.y == prev.y {
				dir := mathext.Sign(curr.x - prev.x)
				for i := prev.x + dir; i != curr.x; i += dir {
					rockSet.Add(Pos2D{i, curr.y})
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

		for curr.y < maxY {
			if part == 2 && curr.y+1 == maxY {
				// Sand hit floor
				rockSet.Add(curr)
				break
			} else if test := (Pos2D{curr.x, curr.y + 1}); !rockSet.Contains(test) {
				curr = test
			} else if test := (Pos2D{curr.x - 1, curr.y + 1}); !rockSet.Contains(test) {
				curr = test
			} else if test := (Pos2D{curr.x + 1, curr.y + 1}); !rockSet.Contains(test) {
				curr = test
			} else {
				// Put sand to rest as rock
				rockSet.Add(curr)
				break
			}
		}

		infiniteFlow := curr.y >= maxY
		sourceBlocked := curr == source
		if infiniteFlow || sourceBlocked {
			if sourceBlocked {
				totalSand++
			}
			break
		}

		totalSand++
	}

	fmt.Println(totalSand)
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
