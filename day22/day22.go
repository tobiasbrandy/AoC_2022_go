package day22

import (
	"errors"
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
)

type Range struct {
	L, R int
}

func (r Range) In(i int) bool {
	return i >= r.L && i < r.R
}

type Dir byte

const (
	RIGHT Dir = '>'
	DOWN  Dir = 'v'
	LEFT  Dir = '<'
	UP    Dir = '^'
)

func (d Dir) Order() int {
	switch d {
	case RIGHT:
		return 0
	case DOWN:
		return 1
	case LEFT:
		return 2
	case UP:
		return 3
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return 0
	}
}

func (d Dir) Sign() int {
	switch d {
	case RIGHT:
		return 1
	case DOWN:
		return 1
	case LEFT:
		return -1
	case UP:
		return -1
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return 0
	}
}

func (d Dir) UpdatePos(p pos.D2, steps int) pos.D2 {
	switch d {
	case RIGHT:
		return pos.New2D(p.X+steps, p.Y)
	case DOWN:
		return pos.New2D(p.X, p.Y+steps)
	case LEFT:
		return pos.New2D(p.X-steps, p.Y)
	case UP:
		return pos.New2D(p.X, p.Y-steps)
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return pos.D2{}
	}
}

func (d Dir) RotateR() Dir {
	switch d {
	case RIGHT:
		return DOWN
	case DOWN:
		return LEFT
	case LEFT:
		return UP
	case UP:
		return RIGHT
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return 0
	}
}

func (d Dir) RotateL() Dir {
	switch d {
	case RIGHT:
		return UP
	case DOWN:
		return RIGHT
	case LEFT:
		return DOWN
	case UP:
		return LEFT
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return 0
	}
}

func (d Dir) String() string {
	return string(d)
}

type OffsetSlice[T any] struct {
	offset int
	slice  []T
}

func (s OffsetSlice[T]) Len() int {
	return s.offset + len(s.slice)
}

func (s OffsetSlice[T]) Offset() int {
	return s.offset
}

func (s OffsetSlice[T]) InBounds(i int) bool {
	return i >= s.offset && i < s.Len()
}

func (s OffsetSlice[T]) Get(i int) T {
	return s.slice[i-s.offset]
}

func (s OffsetSlice[T]) Mod(i int) int {
	return mathext.Mod(i-s.offset, len(s.slice)) + s.offset
}

func Part1(inputPath string) any {
	var passMap []OffsetSlice[byte]

	scanner := fileline.NewScanner(inputPath, errexit.HandleScanError)
	defer scanner.Close()

	scanner.ForEachWhile(func(line string) bool {
		if line == "" {
			return false
		}

		slice := []byte(line)
		for i, b := range slice {
			if b == ' ' {
				continue
			}

			passMap = append(passMap, OffsetSlice[byte]{
				offset: i,
				slice:  slice[i:],
			})
			break
		}

		return true
	})
	passMapH := len(passMap)

	cmdsS, ok := scanner.Read1()
	if !ok {
		errexit.HandleMainError(errors.New("no instruction line after newline"))
	}
	cmds := []byte(cmdsS)
	cmdsC := len(cmds)

	p := pos.New2D(passMap[0].Offset(), 0)
	dir := RIGHT

	i := 0
	for i < cmdsC {
		// Rotation command
		if cmds[i] == 'R' {
			// Right rotation command
			dir = dir.RotateR()
			i++
		} else if cmds[i] == 'L' {
			// Left rotation command
			dir = dir.RotateL()
			i++
		} else {
			// Steps command
			start := i

			for i < cmdsC && cmds[i] >= '0' && cmds[i] <= '9' {
				i++
			}

			if i == start {
				errexit.HandleMainError(fmt.Errorf("invalid map character: %v", string(cmds[start])))
			}

			steps := parse.Int(string(cmds[start:i]))
			for j := 0; j < steps; j++ {
				newP := dir.UpdatePos(p, 1)

				// Check bounds
				if dir == UP || dir == DOWN {
					newP.Y = mathext.Mod(newP.Y, passMapH)
					sign := -dir.Sign()
					if !passMap[newP.Y].InBounds(newP.X) {
						for newP.Y = p.Y; passMap[newP.Y].InBounds(newP.X); newP.Y = mathext.Mod(newP.Y+sign, passMapH) {
						}
						newP.Y = mathext.Mod(newP.Y-sign, passMapH)
					}
				} else {
					newP.X = passMap[newP.Y].Mod(newP.X)
				}

				if passMap[newP.Y].Get(newP.X) == '#' {
					break
				}
				p = newP
			}
		}
	}

	return 1000*(p.Y+1) + 4*(p.X+1) + dir.Order()
}

func Part2(inputPath string) any {
	// Didn't even try
	return 0
}
