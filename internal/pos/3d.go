package pos

import (
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/hashext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"io"
)

type D3 struct {
	X, Y, Z int
}

func New3D(x, y, z int) D3 {
	return D3{x, y, z}
}

func (p D3) Neighbours6() [6]D3 {
	return [...]D3{
		{p.X - 1, p.Y, p.Z},
		{p.X + 1, p.Y, p.Z},
		{p.X, p.Y - 1, p.Z},
		{p.X, p.Y + 1, p.Z},
		{p.X, p.Y, p.Z + 1},
		{p.X, p.Y, p.Z - 1},
	}
}

func (p D3) Neighbours7() [7]D3 {
	return [...]D3{
		{p.X - 1, p.Y, p.Z},
		{p.X + 1, p.Y, p.Z},
		{p.X, p.Y - 1, p.Z},
		p,
		{p.X, p.Y + 1, p.Z},
		{p.X, p.Y, p.Z + 1},
		{p.X, p.Y, p.Z - 1},
	}
}

func (p D3) Neighbours26() [26]D3 {
	return [...]D3{
		{p.X - 1, p.Y - 1, p.Z - 1},
		{p.X - 1, p.Y - 1, p.Z},
		{p.X - 1, p.Y - 1, p.Z + 1},
		{p.X - 1, p.Y, p.Z - 1},
		{p.X - 1, p.Y, p.Z},
		{p.X - 1, p.Y, p.Z + 1},
		{p.X - 1, p.Y + 1, p.Z - 1},
		{p.X - 1, p.Y + 1, p.Z},
		{p.X - 1, p.Y + 1, p.Z + 1},

		{p.X, p.Y - 1, p.Z - 1},
		{p.X, p.Y - 1, p.Z},
		{p.X, p.Y - 1, p.Z + 1},
		{p.X, p.Y, p.Z - 1},
		// {p.X, p.Y, p.Z}, Me
		{p.X, p.Y, p.Z + 1},
		{p.X, p.Y + 1, p.Z - 1},
		{p.X, p.Y + 1, p.Z},
		{p.X, p.Y + 1, p.Z + 1},

		{p.X + 1, p.Y - 1, p.Z - 1},
		{p.X + 1, p.Y - 1, p.Z},
		{p.X + 1, p.Y - 1, p.Z + 1},
		{p.X + 1, p.Y, p.Z - 1},
		{p.X + 1, p.Y, p.Z},
		{p.X + 1, p.Y, p.Z + 1},
		{p.X + 1, p.Y + 1, p.Z - 1},
		{p.X + 1, p.Y + 1, p.Z},
		{p.X + 1, p.Y + 1, p.Z + 1},
	}
}

func (p D3) Neighbours27() [27]D3 {
	return [...]D3{
		{p.X - 1, p.Y - 1, p.Z - 1},
		{p.X - 1, p.Y - 1, p.Z},
		{p.X - 1, p.Y - 1, p.Z + 1},
		{p.X - 1, p.Y, p.Z - 1},
		{p.X - 1, p.Y, p.Z},
		{p.X - 1, p.Y, p.Z + 1},
		{p.X - 1, p.Y + 1, p.Z - 1},
		{p.X - 1, p.Y + 1, p.Z},
		{p.X - 1, p.Y + 1, p.Z + 1},

		{p.X, p.Y - 1, p.Z - 1},
		{p.X, p.Y - 1, p.Z},
		{p.X, p.Y - 1, p.Z + 1},
		{p.X, p.Y, p.Z - 1},
		p,
		{p.X, p.Y, p.Z + 1},
		{p.X, p.Y + 1, p.Z - 1},
		{p.X, p.Y + 1, p.Z},
		{p.X, p.Y + 1, p.Z + 1},

		{p.X + 1, p.Y - 1, p.Z - 1},
		{p.X + 1, p.Y - 1, p.Z},
		{p.X + 1, p.Y - 1, p.Z + 1},
		{p.X + 1, p.Y, p.Z - 1},
		{p.X + 1, p.Y, p.Z},
		{p.X + 1, p.Y, p.Z + 1},
		{p.X + 1, p.Y + 1, p.Z - 1},
		{p.X + 1, p.Y + 1, p.Z},
		{p.X + 1, p.Y + 1, p.Z + 1},
	}
}

func (p D3) Distance1(o D3) int {
	return mathext.IntAbs(p.X-o.X) + mathext.IntAbs(p.Y-o.Y) + mathext.IntAbs(p.Z-o.Z)
}

func (p D3) String() string {
	return fmt.Sprint("(", p.X, ", ", p.Y, ", ", p.Z, ")")
}

func (p D3) Hash(h io.Writer) {
	hashext.HashNum(h, p.X)
	hashext.HashNum(h, p.Y)
	hashext.HashNum(h, p.Z)
}
