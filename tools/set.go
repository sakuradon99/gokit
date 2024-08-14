package tools

import "sort"

type Set[T comparable] interface {
	Put(e T)
	Remove(e T)
	Contains(e T) bool
	ToSlice() []T
	Len() int
}

type UnstableSet[T comparable] struct {
	m map[T]bool
}

func NewUnstableSet[T comparable]() *UnstableSet[T] {
	return &UnstableSet[T]{
		m: make(map[T]bool),
	}
}

func (s *UnstableSet[T]) Put(e T) {
	s.m[e] = true
}

func (s *UnstableSet[T]) Remove(e T) {
	delete(s.m, e)
}

func (s *UnstableSet[T]) Contains(e T) bool {
	_, ok := s.m[e]
	return ok
}

func (s *UnstableSet[T]) ToSlice() []T {
	var arr []T
	for val := range s.m {
		arr = append(arr, val)
	}
	return arr
}

func (s *UnstableSet[T]) Len() int {
	return len(s.m)
}

type StableSet[T comparable] struct {
	m   map[T]uint64
	inc uint64
}

func NewStableSet[T comparable]() *StableSet[T] {
	return &StableSet[T]{
		m: make(map[T]uint64),
	}
}

func (s *StableSet[T]) Put(e T) {
	s.inc++
	s.m[e] = s.inc
}

func (s *StableSet[T]) Remove(e T) {
	delete(s.m, e)
}

func (s *StableSet[T]) Contains(e T) bool {
	_, ok := s.m[e]
	return ok
}

func (s *StableSet[T]) ToSlice() []T {
	arr := make([]T, 0, len(s.m))
	for val := range s.m {
		arr = append(arr, val)
	}
	sort.Slice(arr, func(i, j int) bool {
		return s.m[arr[i]] < s.m[arr[j]]
	})
	return arr
}

func (s *StableSet[T]) Len() int {
	return len(s.m)
}
