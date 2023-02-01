package day24

import (
	"errors"
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/priorityq"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
	"github.com/tobiasbrandy/AoC_2022_go/internal/stringer"
	"strings"
)

type Dir string

const (
	Up    Dir = "^"
	Down  Dir = "v"
	Left  Dir = "<"
	Right Dir = ">"
)

func (d Dir) Forwards(p pos.D2, steps int) pos.D2 {
	switch d {
	case Up:
		return pos.D2{X: p.X, Y: p.Y - steps}
	case Down:
		return pos.D2{X: p.X, Y: p.Y + steps}
	case Left:
		return pos.D2{X: p.X - steps, Y: p.Y}
	case Right:
		return pos.D2{X: p.X + steps, Y: p.Y}
	default:
		errexit.HandleMainError(fmt.Errorf("invalid direction: %v", d))
		return pos.D2{}
	}
}

type Wind struct {
	pos pos.D2
	dir Dir
}

func (w Wind) Pos(steps, width, height int) pos.D2 {
	p := w.dir.Forwards(w.pos, steps)
	p.X, p.Y = mathext.Mod(p.X, width), mathext.Mod(p.Y, height)
	return p
}

type Board struct {
	vWinds, hWinds map[int][]Wind
	width, height  int
	start, end     pos.D2
}

func (b Board) String() string {
	return stringer.String(b)
}

type State struct {
	pos  pos.D2
	time int
	dist int
}

func (s State) Cost() int {
	return s.time + s.dist // Heuristic: distance1 to end
}

func (s State) Less(o State) bool {
	return o.Cost() < s.Cost()
}

func (s State) String() string {
	return stringer.String(s)
}

func parseBoard(inputPath string) Board {
	scanner := fileline.NewScanner(inputPath, errexit.HandleScanError)
	defer scanner.Close()

	board := Board{
		vWinds: make(map[int][]Wind), hWinds: make(map[int][]Wind),
		width: 0, height: 0,
		start: pos.D2{}, end: pos.D2{},
	}

	firstLine, _ := scanner.Read1()
	board.width = len(firstLine) - 2
	board.start = pos.New2D(strings.IndexByte(firstLine, '.')-1, -1)

	y := 0
	scanner.ForEachWhile(func(line string) bool {
		if strings.Count(line[:3], "#") > 1 {
			// End line
			board.height = y
			board.end = pos.New2D(strings.IndexByte(line, '.')-1, board.height)
			return false
		}

		for i, c := range line[1:] {
			dir := Dir(c)
			switch dir {
			case Up, Down:
				board.vWinds[i] = append(board.vWinds[i], Wind{
					pos: pos.New2D(i, y),
					dir: dir,
				})
			case Right, Left:
				board.hWinds[y] = append(board.hWinds[y], Wind{
					pos: pos.New2D(i, y),
					dir: dir,
				})
			}
		}

		y++
		return true
	})

	return board
}

func validState(board *Board, state State) bool {
	p, time := state.pos, state.time
	if p.X < 0 || p.X >= board.width {
		return false
	}
	if p.Y < 0 || p.Y >= board.height {
		return false
	}
	for _, wind := range board.vWinds[p.X] {
		if wind.Pos(time, board.width, board.height) == p {
			return false
		}
	}
	for _, wind := range board.hWinds[p.Y] {
		if wind.Pos(time, board.width, board.height) == p {
			return false
		}
	}
	return true
}

func bestTravelTime(board *Board, initTime int) int {
	// We target the space before reaching the end if it is located in a border
	end := board.end
	if end.Y == -1 {
		end.Y++
	} else if end.Y == board.height {
		end.Y--
	}

	states := priorityq.New[State]()
	visited := set.Set[State]{}

	states.Push(State{
		pos:  board.start,
		time: initTime,
		dist: board.start.Distance1(end),
	})

	for states.Len() > 0 {
		state := states.Pop() // Dijkstra
		p := state.pos
		time := state.time + 1

		if p == end {
			// End state - We know it's best time because Dijkstra
			return time
		}

		// Cardinal directions and wait
		for _, newPos := range p.Neighbours5() {
			newState := State{
				pos:  newPos,
				time: time,
				dist: newPos.Distance1(end),
			}

			if !visited.Contains(newState) && (newPos == board.start || validState(board, newState)) {
				visited.Add(newState)
				states.Push(newState)
			}
		}
	}

	errexit.HandleMainError(errors.New("couldn't reach end"))
	return 0
}

func Part1(inputPath string) any {
	board := parseBoard(inputPath)
	return bestTravelTime(&board, 0)
}

func Part2(inputPath string) any {
	board := parseBoard(inputPath)
	time := bestTravelTime(&board, 0)

	board.start, board.end = board.end, board.start
	time = bestTravelTime(&board, time)

	board.start, board.end = board.end, board.start
	time = bestTravelTime(&board, time)

	return time
}
