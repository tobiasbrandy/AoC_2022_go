package internal

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](capacity int) Set[T] {
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
