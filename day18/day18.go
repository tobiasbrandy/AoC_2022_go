package day18

import (
	"github.com/gammazero/deque"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
	"math"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
)

func parseVertex(data string) pos.D3 {
	vs := strings.Split(data, ",")
	return pos.New3D(parse.Int(vs[0]), parse.Int(vs[1]), parse.Int(vs[2]))
}

func Part1(inputPath string) any {
	vertices := set.Set[pos.D3]{}

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		vertices.Add(parseVertex(line))
	})

	perimeter := 0
	for v := range vertices {
		for _, neighbour := range v.Neighbours6() {
			if !vertices.Contains(neighbour) {
				perimeter++
			}
		}
	}

	return perimeter
}

func Part2(inputPath string) any {
	vertices := set.Set[pos.D3]{}
	start := pos.New3D(math.MaxInt, math.MaxInt, math.MaxInt)

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		v := parseVertex(line)
		vertices.Add(v)

		// We guarantee start to be in the left perimeter
		if start.X > v.X {
			start = v
		}
	})

	// We start at the left empty space of start
	start.X--

	tasks := deque.New[pos.D3](vertices.Len())
	tasks.PushBack(start)

	touched := set.New[pos.D3](vertices.Len())
	touched.Add(start)

	perimeter := 0
	for tasks.Len() > 0 {
		v := tasks.PopFront()
		nToPush := make([]pos.D3, 0, 6)
		isPerimeter := false

		for _, neighbour := range v.Neighbours6() {
			if vertices.Contains(neighbour) {
				perimeter++
				isPerimeter = true
			} else {
				if !touched.Contains(neighbour) {
					nToPush = append(nToPush, neighbour)
				}
			}
		}

		if !isPerimeter {
			// Long perimeter check: 26-connected
			for _, neighbour := range v.Neighbours26() {
				if vertices.Contains(neighbour) {
					isPerimeter = true
					break
				}
			}
		}

		if isPerimeter {
			for _, neighbour := range nToPush {
				tasks.PushBack(neighbour)
				touched.Add(neighbour)
			}
		}
	}

	return perimeter
}
