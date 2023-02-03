package set

import (
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/hashext"
	"io"
	"strings"
)

type Set[T comparable] map[T]struct{}

func New[T comparable](capacity int) Set[T] {
	return make(map[T]struct{}, capacity)
}

func (set Set[T]) Add(elem T) {
	set[elem] = struct{}{}
}

func (set Set[T]) AddAll(elems []T) {
	for _, elem := range elems {
		set.Add(elem)
	}
}

func (set Set[T]) Remove(elem T) {
	delete(set, elem)
}

func (set Set[T]) Contains(elem T) bool {
	_, ok := set[elem]
	return ok
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) Copy() Set[T] {
	ret := New[T](set.Len())

	for e := range set {
		ret.Add(e)
	}

	return ret
}

func (set Set[T]) Diff(o Set[T]) Set[T] {
	ret := New[T](set.Len())

	for e := range set {
		if !o.Contains(e) {
			ret.Add(e)
		}
	}

	return ret
}

func (set Set[T]) Disjoint(o Set[T]) bool {
	return set.Diff(o).Len() == set.Len()
}

func (set Set[T]) Hash(h io.Writer) {
	for v := range set {
		hashext.HashAny(h, v)
	}
}

func (set Set[T]) String() string {
	var sb strings.Builder

	sb.WriteString("[ ")
	for elem := range set {
		sb.WriteString(fmt.Sprint(elem))
		sb.WriteRune(' ')
	}
	sb.WriteRune(']')

	return sb.String()
}
