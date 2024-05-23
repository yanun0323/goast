package goast

type set[T comparable] map[T]struct{}

func (s *set[T]) Contain(key T) bool {
	if s == nil || *s == nil {
		return false
	}
	_, ok := (*s)[key]
	return ok
}

func (s *set[T]) Insert(key T) {
	if s == nil {
		s = &set[T]{}
	}

	if *s == nil {
		*s = set[T]{}
	}

	(*s)[key] = struct{}{}
}
