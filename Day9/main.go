package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"strconv"
	"flag"
	"fmt"
)

type Pos2D struct {
	x, y int
}

func solve(filePath string, part int) {
	var nodeCount int
	if part == 1 {
		nodeCount = 2
	} else { // part == 2
		nodeCount = 10
	}

	// Nodes implicitly initialized on {0, 0}
	nodes := make([]Pos2D, nodeCount)
	head := &nodes[0]
	tail := &nodes[nodeCount-1]

	tailPosCache := internal.Set[Pos2D]{}
	tailPosCache.Add(*tail)

	internal.ForEachFileLine(filePath, internal.HandleScanError, func(line string) {
		// Parse
		cmd := []rune(line[0:1])[0]
		count, err := strconv.Atoi(line[2:])
		if err != nil {
			internal.HandleMainError(err)
		}

		// Move head
		for i := 0; i < count; i++ {
			switch cmd {
			case 'R':
				head.x++
			case 'L':
				head.x--
			case 'U':
				head.y++
			case 'D':
				head.y--
			default:
				internal.HandleMainError(fmt.Errorf("invalid move direction %v", cmd))
			}

			for i := 1; i < nodeCount; i++ {
				prev := &nodes[i-1]
				curr := &nodes[i]

				dx := prev.x - curr.x
				dy := prev.y - curr.y
				mdx := internal.IntAbs(dx)
				mdy := internal.IntAbs(dy)
	
				if mdx > 1 {
					curr.x += internal.Sign(dx)
					if mdy > 0 {
						curr.y += internal.Sign(dy)
					}
				} else if mdy > 1 {
					curr.y += internal.Sign(dy)
					if mdx > 0 {
						curr.x += internal.Sign(dx)
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
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}

	solve(*inputPath, *part)
}
