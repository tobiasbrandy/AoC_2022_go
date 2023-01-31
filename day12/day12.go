package day12

import (
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"

	"github.com/gammazero/deque"
)

func Solve(inputPath string, part int) any {
	var heights [][]int

	starts := make([]pos.D2, 1)
	var end pos.D2

	row := 0
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		heights = append(heights, make([]int, len(line)))

		for col, level := range line {
			if level == 'S' {
				heights[row][col] = int('a')
				starts[0] = pos.New2D(row, col)
			} else if level == 'E' {
				heights[row][col] = int('z')
				end = pos.New2D(row, col)
			} else if part == 2 && level == 'a' {
				starts = append(starts, pos.New2D(row, col))
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
		visited := make(map[pos.D2]int, width*height) // How many steps to get there
		visited[start] = 0

		processQueue := &deque.Deque[pos.D2]{} // DFS
		processQueue.PushBack(start)

		for processQueue.Len() > 0 {
			p := processQueue.PopFront()
			h := heights[p.X][p.Y]
			steps := visited[p] + 1

			for _, n := range p.Neighbours4() {
				if n.X >= 0 && n.X < width && n.Y >= 0 && n.Y < height && heights[n.X][n.Y] <= h+1 { // Inbounds and height
					nSteps, ok := visited[n]
					if !ok || nSteps > steps {
						visited[n] = steps
						processQueue.PushBack(n)
					}
				}
			}
		}

		l, ok := visited[end]
		if ok && l < shortestLen {
			shortestLen = l
		}
	}

	return shortestLen
}
