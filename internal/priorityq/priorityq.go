package priorityq

import (
	"container/heap"
	"fmt"
	"github.com/gammazero/deque"
	"strings"
)

type Ordered[T any] interface {
	Less(o T) bool
}

type dequeHeap[T Ordered[T]] deque.Deque[T]

func (pq *dequeHeap[T]) Len() int {
	return (*deque.Deque[T])(pq).Len()
}

func (pq *dequeHeap[T]) Less(i, j int) bool {
	dq := (*deque.Deque[T])(pq)
	return dq.At(j).Less(dq.At(i))
}

func (pq *dequeHeap[T]) Swap(i, j int) {
	dq := (*deque.Deque[T])(pq)
	ei, ej := dq.At(i), dq.At(j)
	dq.Set(i, ej)
	dq.Set(j, ei)
}

func (pq *dequeHeap[T]) Push(x any) {
	(*deque.Deque[T])(pq).PushBack(x.(T))
}

func (pq *dequeHeap[T]) Pop() any {
	return (*deque.Deque[T])(pq).PopBack()
}

type PriorityQueue[T Ordered[T]] dequeHeap[T]

func New[T Ordered[T]](size ...int) *PriorityQueue[T] {
	return (*PriorityQueue[T])(deque.New[T](size...))
}

func (pq *PriorityQueue[T]) Cap() int {
	return (*deque.Deque[T])(pq).Cap()
}

func (pq *PriorityQueue[T]) Len() int {
	return (*deque.Deque[T])(pq).Len()
}

func (pq *PriorityQueue[T]) Push(elem T) {
	heap.Push((*dequeHeap[T])(pq), elem)
}

func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop((*dequeHeap[T])(pq)).(T)
}

func (pq *PriorityQueue[T]) Peek() T {
	return (*deque.Deque[T])(pq).Front()
}

func (pq *PriorityQueue[T]) At(i int) T {
	return (*deque.Deque[T])(pq).At(i)
}

func (pq *PriorityQueue[T]) Set(i int, item T) {
	(*deque.Deque[T])(pq).Set(i, item)
	heap.Fix((*dequeHeap[T])(pq), i)
}

func (pq *PriorityQueue[T]) Clear() {
	(*deque.Deque[T])(pq).Clear()
}

func (pq *PriorityQueue[T]) Index(f func(T) bool) int {
	return (*deque.Deque[T])(pq).Index(f)
}

func (pq *PriorityQueue[T]) Remove(at int) T {
	return heap.Remove((*dequeHeap[T])(pq), at).(T)
}

func (pq *PriorityQueue[T]) SetMinCapacity(minCapacityExp uint) {
	(*deque.Deque[T])(pq).SetMinCapacity(minCapacityExp)
}

func (pq *PriorityQueue[T]) String() string {
	l := pq.Len()
	if l == 0 {
		return "[]"
	}

	tmp := make([]T, l, l)
	builder := strings.Builder{}

	builder.WriteString("[")
	e := pq.Pop()
	tmp[0] = e
	builder.WriteString(fmt.Sprint(e))

	for i := 1; i < l; i++ {
		e := pq.Pop()
		tmp[i] = e
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprint(e))
	}
	builder.WriteString("]")

	for i := 0; i < l; i++ {
		pq.Push(tmp[i])
	}

	return builder.String()
}
