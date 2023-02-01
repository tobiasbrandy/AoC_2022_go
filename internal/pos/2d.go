package pos

import (
	"fmt"
	"io"

	"github.com/tobiasbrandy/AoC_2022_go/internal/hashext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
)

type D2 struct {
	X, Y int
}

func New2D(x, y int) D2 {
	return D2{x, y}
}

func (p D2) Neighbours4() [4]D2 {
	return [...]D2{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func (p D2) Neighbours5() [5]D2 {
	return [...]D2{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		p,
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func (p D2) Neighbours8() [8]D2 {
	return [...]D2{
		{p.X - 1, p.Y - 1},
		{p.X - 1, p.Y},
		{p.X - 1, p.Y + 1},

		{p.X, p.Y - 1},
		// {p.X, p.Y}, Me
		{p.X, p.Y + 1},

		{p.X + 1, p.Y - 1},
		{p.X + 1, p.Y},
		{p.X + 1, p.Y + 1},
	}
}

func (p D2) Neighbours9() [9]D2 {
	return [...]D2{
		{p.X - 1, p.Y - 1},
		{p.X - 1, p.Y},
		{p.X - 1, p.Y + 1},

		{p.X, p.Y - 1},
		p,
		{p.X, p.Y + 1},

		{p.X + 1, p.Y - 1},
		{p.X + 1, p.Y},
		{p.X + 1, p.Y + 1},
	}
}

func (p D2) Distance1(o D2) int {
	return mathext.IntAbs(p.X-o.X) + mathext.IntAbs(p.Y-o.Y)
}

func (p D2) String() string {
	return fmt.Sprint("(", p.X, ", ", p.Y, ")")
}

func (p D2) Hash(h io.Writer) {
	hashext.HashNum(h, p.X)
	hashext.HashNum(h, p.Y)
}
