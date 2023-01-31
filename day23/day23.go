package day23

import (
	"fmt"
	"github.com/gammazero/deque"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
	"math"
)

type Dir string

const (
	Undefined Dir = ""
	North     Dir = "North"
	South     Dir = "South"
	West      Dir = "West"
	East      Dir = "East"
	DirC      int = 4
)

func Dirs() [DirC]Dir {
	return [...]Dir{North, South, West, East}
}

func (d Dir) Neighbours3(p pos.D2) [3]pos.D2 {
	switch d {
	case North:
		return [...]pos.D2{
			{X: p.X - 1, Y: p.Y - 1},
			{X: p.X, Y: p.Y - 1},
			{X: p.X + 1, Y: p.Y - 1},
		}
	case South:
		return [...]pos.D2{
			{X: p.X - 1, Y: p.Y + 1},
			{X: p.X, Y: p.Y + 1},
			{X: p.X + 1, Y: p.Y + 1},
		}
	case East:
		return [...]pos.D2{
			{X: p.X + 1, Y: p.Y - 1},
			{X: p.X + 1, Y: p.Y},
			{X: p.X + 1, Y: p.Y + 1},
		}
	case West:
		return [...]pos.D2{
			{X: p.X - 1, Y: p.Y - 1},
			{X: p.X - 1, Y: p.Y},
			{X: p.X - 1, Y: p.Y + 1},
		}
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return [3]pos.D2{{}, {}, {}}
	}
}

func (d Dir) Forwards(p pos.D2) pos.D2 {
	switch d {
	case North:
		return pos.D2{X: p.X, Y: p.Y - 1}
	case South:
		return pos.D2{X: p.X, Y: p.Y + 1}
	case East:
		return pos.D2{X: p.X + 1, Y: p.Y}
	case West:
		return pos.D2{X: p.X - 1, Y: p.Y}
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return pos.D2{}
	}
}

func (d Dir) Backwards(p pos.D2) pos.D2 {
	switch d {
	case North:
		return pos.D2{X: p.X, Y: p.Y + 1}
	case South:
		return pos.D2{X: p.X, Y: p.Y - 1}
	case East:
		return pos.D2{X: p.X - 1, Y: p.Y}
	case West:
		return pos.D2{X: p.X + 1, Y: p.Y}
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return pos.D2{}
	}
}

func parseInitPositions(inputPath string) set.Set[pos.D2] {
	positions := set.Set[pos.D2]{}

	y := 0
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		for x, c := range line {
			if c == '#' {
				positions.Add(pos.New2D(x, y))
			}
		}
		y++
	})

	return positions
}

func posHasNeighbour8(p pos.D2, positions set.Set[pos.D2]) bool {
	for _, n := range p.Neighbours8() {
		if positions.Contains(n) {
			return true
		}
	}
	return false
}

func posHasDirNeighbour3(p pos.D2, dir Dir, positions set.Set[pos.D2]) bool {
	for _, n := range dir.Neighbours3(p) {
		if positions.Contains(n) {
			return true
		}
	}
	return false
}

func availableDir(p pos.D2, dirs deque.Deque[Dir], positions set.Set[pos.D2]) (Dir, bool) {
	for di := 0; di < DirC; di++ {
		dir := dirs.At(di)
		if !posHasDirNeighbour3(p, dir, positions) {
			return dir, true
		}
	}
	return Undefined, false
}

func updatePositions(positions set.Set[pos.D2], dirs deque.Deque[Dir]) (newPositions set.Set[pos.D2], hasChanged bool) {
	positionsC := positions.Len()
	newPositions = set.New[pos.D2](positionsC)
	collisions := set.New[pos.D2](positionsC / 2)
	newToOldPos := make(map[pos.D2]pos.D2, positionsC)
	changed := 0

	for p := range positions {
		if posHasNeighbour8(p, positions) {
			if dir, available := availableDir(p, dirs, positions); available {
				newPos := dir.Forwards(p)
				if collisions.Contains(newPos) {
					newPositions.Add(p)
				} else if newPositions.Contains(newPos) {
					collisions.Add(newPos)
					newPositions.Add(p)

					// Remove collision from newPositions
					newPositions.Remove(newPos)
					newPositions.Add(newToOldPos[newPos])
					changed--
				} else {
					newPositions.Add(newPos)
					newToOldPos[newPos] = p
					changed++
				}
			} else {
				newPositions.Add(p)
			}
		} else {
			newPositions.Add(p)
		}
	}

	hasChanged = changed > 0
	return newPositions, hasChanged
}

func Part1(inputPath string) any {
	const rounds int = 10

	positions := parseInitPositions(inputPath)

	dirs := deque.Deque[Dir]{}
	for _, dir := range Dirs() {
		dirs.PushBack(dir)
	}

	for round := 0; round < rounds; round++ {
		var hasChanged bool
		positions, hasChanged = updatePositions(positions, dirs)
		if !hasChanged {
			// Fixed point! -> We can end early
			break
		}

		dirs.Rotate(1)
	}

	// Find bounding box
	minY, minX := math.MaxInt, math.MaxInt
	maxY, maxX := math.MinInt, math.MinInt
	for p := range positions {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	boundingBoxArea := (maxX - minX + 1) * (maxY - minY + 1)
	emptyPos := boundingBoxArea - positions.Len()

	return emptyPos
}

func Part2(inputPath string) any {
	positions := parseInitPositions(inputPath)

	dirs := deque.Deque[Dir]{}
	for _, dir := range Dirs() {
		dirs.PushBack(dir)
	}

	round := 1
	for {
		var hasChanged bool
		positions, hasChanged = updatePositions(positions, dirs)
		if !hasChanged {
			// Fixed point! -> We can end early
			break
		}

		dirs.Rotate(1)
		round++
	}

	return round
}
