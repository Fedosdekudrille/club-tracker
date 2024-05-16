package queue

type node[T comparable] struct {
	next  *node[T]
	value T
}

type Queue[T comparable] struct {
	head *node[T]
	tail *node[T]
	len  int
}

func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{
		head: nil,
		tail: nil,
		len:  0,
	}
}

func (q *Queue[T]) Push(v T) {
	node := &node[T]{value: v}
	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.len++
}

func (q *Queue[T]) Pop() T {
	if q.head == nil {
		return *new(T)
	}
	v := q.head.value
	q.head = q.head.next
	q.len--
	return v
}

func (q *Queue[T]) Peek() T {
	if q.head == nil {
		return *new(T)
	}
	return q.head.value
}

func (q *Queue[T]) IsEmpty() bool {
	return q.head == nil
}

func (q *Queue[T]) Len() int {
	return q.len
}

func (q *Queue[T]) Clear() {
	q.head = nil
	q.tail = nil
	q.len = 0
}
