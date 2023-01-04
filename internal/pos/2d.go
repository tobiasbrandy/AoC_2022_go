package pos

import "fmt"

type D2 struct {
	X, Y int
}

func New2D(x, y int) D2 {
	return D2{x, y}
}

func (p D2) Neighbours4() [4]D2 {
	return [4]D2{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func (p D2) String() string {
	return fmt.Sprint("(", p.X, ", ", p.Y, ")")
}
