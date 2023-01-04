package main

import (
	"flag"
	"fmt"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"

	"github.com/gammazero/deque"
)

type Pos2D struct {
	x, y int
}

func (p Pos2D) Neighbours4() [4]Pos2D {
	return [4]Pos2D{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
}

func solve(filePath string, part int) {
	var heights [][]int

	starts := make([]Pos2D, 1)
	var end Pos2D

	row := 0
	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
		heights = append(heights, make([]int, len(line)))

		for col, level := range line {
			if level == 'S' {
				heights[row][col] = int('a')
				starts[0] = Pos2D{row, col}
			} else if level == 'E' {
				heights[row][col] = int('z')
				end = Pos2D{row, col}
			} else if part == 2 && level == 'a' {
				starts = append(starts, Pos2D{row, col})
				heights[row][col] = int(level)
			} else {
				heights[row][col] = int(level)
			}
		}

		row++
	})

	width := len(heights)
	height := len(heights[0])

	shortestLen := width * height

	for _, start := range starts {
		visited := make(map[Pos2D]int, width*height) // How many steps to get there
		visited[start] = 0

		processQueue := &deque.Deque[Pos2D]{} // DFS
		processQueue.PushBack(start)

		for processQueue.Len() > 0 {
			pos := processQueue.PopFront()
			h := heights[pos.x][pos.y]
			steps := visited[pos] + 1

			for _, n := range pos.Neighbours4() {
				if n.x >= 0 && n.x < width && n.y >= 0 && n.y < height && heights[n.x][n.y] <= h+1 { // Inbounds and height
					nSteps, ok := visited[n]
					if !ok || nSteps > steps {
						visited[n] = steps
						processQueue.PushBack(n)
					}
				}
			}
		}

		len, ok := visited[end]
		if ok && len < shortestLen {
			shortestLen = len
		}
	}

	fmt.Println(shortestLen)
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
