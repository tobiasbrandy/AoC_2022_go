package main

import (
	"flag"
	"fmt"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

func solve(filePath string, part int) {
	var nodeCount int
	if part == 1 {
		nodeCount = 2
	} else { // part == 2
		nodeCount = 10
	}

	// Nodes implicitly initialized on {0, 0}
	nodes := make([]pos.D2, nodeCount)
	head := &nodes[0]
	tail := &nodes[nodeCount-1]

	tailPosCache := set.Set[pos.D2]{}
	tailPosCache.Add(*tail)

	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
		// Parse
		cmd := []rune(line[0:1])[0]
		count := parse.Int(line[2:])

		// Move head
		for i := 0; i < count; i++ {
			switch cmd {
			case 'R':
				head.X++
			case 'L':
				head.X--
			case 'U':
				head.Y++
			case 'D':
				head.Y--
			default:
				errexit.HandleMainError(fmt.Errorf("invalid move direction %v", cmd))
			}

			for i := 1; i < nodeCount; i++ {
				prev := &nodes[i-1]
				curr := &nodes[i]

				dx := prev.X - curr.X
				dy := prev.Y - curr.Y
				mdx := mathext.IntAbs(dx)
				mdy := mathext.IntAbs(dy)

				if mdx > 1 {
					curr.X += mathext.Sign(dx)
					if mdy > 0 {
						curr.Y += mathext.Sign(dy)
					}
				} else if mdy > 1 {
					curr.Y += mathext.Sign(dy)
					if mdx > 0 {
						curr.X += mathext.Sign(dx)
					}
				}
			}

			tailPosCache.Add(*tail)
		}
	})

	fmt.Println(tailPosCache.Len())
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
