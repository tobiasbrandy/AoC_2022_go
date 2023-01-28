package linkedlist

import (
	"fmt"
	"strings"
)

type CircularDouble[T any] struct {
	Data       T
	Prev, Next *CircularDouble[T]
}

func NewCircularDouble[T any](data T) *CircularDouble[T] {
	node := &CircularDouble[T]{
		Data: data,
	}
	node.Next = node
	node.Prev = node
	return node
}

func (l *CircularDouble[T]) Get(n int) *CircularDouble[T] {
	if n == 0 || l == l.Next {
		return l
	} else if n > 0 {
		curr := l.Next
		for i := 1; i < n; i++ {
			curr = curr.Next
		}
		return curr
	} else {
		n = -n
		curr := l.Prev
		for i := 1; i < n; i++ {
			curr = curr.Prev
		}
		return curr
	}
}

func (l *CircularDouble[T]) Append(data T) *CircularDouble[T] {
	node := &CircularDouble[T]{
		Data: data,
	}
	l.AppendNode(node)
	return node
}

func (l *CircularDouble[T]) AppendNode(node *CircularDouble[T]) {
	next := l.Next
	next.Prev = node
	l.Next = node
	node.Next = next
	node.Prev = l
}

func (l *CircularDouble[T]) Remove() {
	l.Prev.Next = l.Next
	l.Next.Prev = l.Prev
	l.Next = l
	l.Prev = l
}

func (l *CircularDouble[T]) Move(n int) {
	if n == 0 || l.Next == l {
		return
	} else if n > 0 {
		curr := l.Next
		l.Remove()
		for i := 1; i < n; i++ {
			curr = curr.Next
		}
		curr.AppendNode(l)
	} else {
		n = -n
		curr := l.Prev
		l.Remove()
		for i := 1; i < n; i++ {
			curr = curr.Prev
		}
		curr.Prev.AppendNode(l)
	}
}

func (l *CircularDouble[T]) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")

	builder.WriteString(fmt.Sprint(l.Data))
	curr := l.Next

	for curr != l {
		builder.WriteString(fmt.Sprint(" ", curr.Data))
		curr = curr.Next
	}

	builder.WriteString("]")
	return builder.String()
}
