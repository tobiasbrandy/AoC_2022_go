package mathext

import (
	"golang.org/x/exp/constraints"
	"math"
)

func IntAbs[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Sign[T constraints.Signed | constraints.Unsigned | constraints.Float](x T) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func IntMax(is ...int) int {
	max := math.MinInt
	for _, i := range is {
		if i > max {
			max = i
		}
	}
	return max
}

func IntPow(x int, exp uint) int {
	ret := 1

	for {
		if exp&1 != 0 {
			ret *= x
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		x *= x
	}

	return ret
}

// Mod is like %, but always returns a positive number or zero
func Mod(x, m int) int {
	return (x%m + m) % m
}
