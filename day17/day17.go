package day17

import (
	"fmt"
	"hash/maphash"
	"io"
	"os"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/hashext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"

	"github.com/gammazero/deque"
)

type Rock struct {
	blocks []pos.D2
	right  []pos.D2
	left   []pos.D2
	down   []pos.D2
}

var rocks = [...]Rock{

	// ####
	{
		blocks: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
			pos.New2D(2, 0),
			pos.New2D(3, 0),
		},
		right: []pos.D2{
			pos.New2D(3, 0),
		},
		left: []pos.D2{
			pos.New2D(0, 0),
		},
		down: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
			pos.New2D(2, 0),
			pos.New2D(3, 0),
		},
	},

	// .#.
	// ###
	// .#.
	{
		blocks: []pos.D2{
			pos.New2D(1, 0),
			pos.New2D(0, 1),
			pos.New2D(1, 1),
			pos.New2D(2, 1),
			pos.New2D(1, 2),
		},
		right: []pos.D2{
			pos.New2D(1, 0),
			pos.New2D(2, 1),
			pos.New2D(1, 2),
		},
		left: []pos.D2{
			pos.New2D(1, 0),
			pos.New2D(0, 1),
			pos.New2D(1, 2),
		},
		down: []pos.D2{
			pos.New2D(1, 0),
			pos.New2D(0, 1),
			pos.New2D(2, 1),
		},
	},

	// ..#
	// ..#
	// ###
	{
		blocks: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
			pos.New2D(2, 0),
			pos.New2D(2, 1),
			pos.New2D(2, 2),
		},
		right: []pos.D2{
			pos.New2D(2, 0),
			pos.New2D(2, 1),
			pos.New2D(2, 2),
		},
		left: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(2, 1),
			pos.New2D(2, 2),
		},
		down: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
			pos.New2D(2, 0),
		},
	},

	// #
	// #
	// #
	// #
	{
		blocks: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(0, 1),
			pos.New2D(0, 2),
			pos.New2D(0, 3),
		},
		right: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(0, 1),
			pos.New2D(0, 2),
			pos.New2D(0, 3),
		},
		left: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(0, 1),
			pos.New2D(0, 2),
			pos.New2D(0, 3),
		},
		down: []pos.D2{
			pos.New2D(0, 0),
		},
	},

	// ##
	// ##
	{
		blocks: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
			pos.New2D(0, 1),
			pos.New2D(1, 1),
		},
		right: []pos.D2{
			pos.New2D(1, 0),
			pos.New2D(1, 1),
		},
		left: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(0, 1),
		},
		down: []pos.D2{
			pos.New2D(0, 0),
			pos.New2D(1, 0),
		},
	},
}

type BoardSurface struct {
	board  *deque.Deque[[]bool]
	bottom int
	width  int
}

func NewBoardSurface(width, bottom, initLen int) *BoardSurface {
	if width < 0 {
		errexit.HandleMainError(fmt.Errorf("width must be positive: %v", width))
	}

	return &BoardSurface{
		board:  deque.New[[]bool](initLen),
		bottom: bottom,
		width:  width,
	}
}

func (bs *BoardSurface) Contains(x, y int) bool {
	y -= bs.bottom

	if x < 0 || x >= bs.width || y < 0 {
		return true // Out of bounds
	}

	if y >= bs.board.Len() {
		return false // Empty
	}

	return bs.board.At(y)[x]
}

func (bs *BoardSurface) UpdateBottom() {
	top := make([]int, bs.width)

	l := bs.board.Len()
	for y := 0; y < l; y++ {
		row := bs.board.At(y)

		full := true
		for x := 0; x < bs.width; x++ {
			if row[x] {
				top[x] = y + 1
			}
			if top[x] == 0 {
				full = false
			}
		}

		if full {
			min := l + 1
			for i := 0; i < bs.width; i++ {
				if min > top[i] {
					min = top[i]
				}
			}

			// Update bottom
			bs.IncBottom(min)
			return
		}
	}
}

func (bs *BoardSurface) IncBottom(n int) {
	if n < 1 {
		return
	}

	bs.bottom += n
	for i := 0; i < n; i++ {
		bs.board.PopFront()
	}
}

func (bs *BoardSurface) Add(x, y int) {
	if x < 0 || x >= bs.width {
		return // Out of bounds
	}

	y -= bs.bottom

	if y < 0 {
		return // Not on surface
	}

	l := bs.board.Len()
	if y >= l {
		for i := l; i <= y; i++ {
			bs.board.PushBack(make([]bool, bs.width))
		}
	}

	bs.board.At(y)[x] = true
}

func (bs *BoardSurface) Hash(h io.Writer) {
	// Hash independent of bottom position
	hashext.HashNum(h, bs.width)

	l := bs.board.Len()
	for i := 0; i < l; i++ {
		hashext.HashNumArr(h, bs.board.At(i))
	}
}

func (bs *BoardSurface) String() string {
	var builder strings.Builder

	for i := bs.board.Len() - 1; i >= 0; i-- {
		row := bs.board.At(i)

		for i := 0; i < bs.width; i++ {
			if row[i] {
				builder.WriteRune('#')
			} else {
				builder.WriteRune('.')
			}
		}
		builder.WriteRune('\n')
	}

	builder.WriteString(fmt.Sprint("Bottom:", bs.bottom))

	return builder.String()
}

func simStateHash(h *maphash.Hash, bs *BoardSurface, cmd, rock int) uint64 {
	bs.Hash(h)
	hashext.HashNum(h, cmd)
	hashext.HashNum(h, rock)

	ret := h.Sum64()
	h.Reset()
	return ret
}

func Solve(inputPath string, part int) any {
	const (
		sepY int = 3
		sepX int = 2

		boardWidth int = 7
	)

	var totalRocks int
	if part == 1 {
		totalRocks = 2022
	} else { // part == 2
		totalRocks = 1_000_000_000_000
	}

	file, err := os.Open(inputPath)
	if err != nil {
		errexit.HandleScanError(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			errexit.HandleScanError(err)
		}
	}(file)

	cmds, err := io.ReadAll(file)
	if err != nil {
		errexit.HandleScanError(err)
	}

	rocksC, cmdsC := len(rocks), len(cmds)
	maxY := -1

	var board = NewBoardSurface(boardWidth, 0, 1+sepY)

	var h maphash.Hash
	cycleCache := make(map[uint64]int)
	rockHeights := make([]int, 0)

	for r, c := 0, -1; r < totalRocks; r++ {
		// Cycle analysis
		simStateH := simStateHash(&h, board, c%cmdsC, r%rocksC)
		if cycleStart, ok := cycleCache[simStateH]; ok {
			// Cycle found -> We can extrapolate the rest of the simulation
			cycleLen := r - cycleStart
			rocksLeft := totalRocks - r
			cycles := rocksLeft / cycleLen
			cycleHeight := maxY - rockHeights[cycleStart]
			cycleRemainder := rocksLeft % cycleLen
			remainderHeight := rockHeights[cycleStart+cycleRemainder] - rockHeights[cycleStart] // cycleRemainder height gained from cycleStart
			maxY += cycles*cycleHeight + remainderHeight
			break // End simulation
		} else {
			cycleCache[simStateH] = r
			rockHeights = append(rockHeights, maxY)
		}

		rock := rocks[r%rocksC]
		curr := pos.New2D(sepX, maxY+1+sepY)

		for {
			c++
			cmd := cmds[c%cmdsC]

			if cmd == '<' {
				valid := true
				for _, b := range rock.left {
					newX := curr.X + b.X - 1
					if board.Contains(newX, curr.Y+b.Y) {
						valid = false
						break
					}
				}
				if valid {
					curr.X--
				}
			} else if cmd == '>' {
				valid := true
				for _, b := range rock.right {
					newX := curr.X + b.X + 1
					if board.Contains(newX, curr.Y+b.Y) {
						valid = false
						break
					}
				}
				if valid {
					curr.X++
				}
			}

			rest := false
			for _, b := range rock.down {
				newY := curr.Y + b.Y - 1
				if board.Contains(curr.X+b.X, newY) {
					rest = true
					break
				}
			}
			if rest {
				for _, b := range rock.blocks {
					newX, newY := curr.X+b.X, curr.Y+b.Y
					if newY > maxY {
						maxY = newY
					}
					board.Add(newX, newY)
				}
				board.UpdateBottom()
				break
			}

			curr.Y--
		}
	}

	height := maxY + 1
	return height
}
