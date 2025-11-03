package utils

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]struct{}),
	}

	s.Add(items...)

	return s
}

func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *Set[T]) Size() int {
	return len(s.m)
}
