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
