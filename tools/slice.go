package tools

import (
	"math/rand"
	"time"
)

func NewSlice[T any](size ...int) []T {
	s := 0
	if len(size) > 0 {
		s = size[len(size)-1]
	}
	return make([]T, s)
}

func SliceContains[T comparable](arr []T, val T) bool {
	for _, t := range arr {
		if val == t {
			return true
		}
	}
	return false
}

func SliceContainsF[T any](arr []T, val T, f func(src, dst T) bool) bool {
	for _, t := range arr {
		if f(t, val) {
			return true
		}
	}
	return false
}

func SliceGetOrDefault[T any](arr []T, i int, def T) T {
	if len(arr)-1 < i {
		return def
	}
	return arr[i]
}

func SliceConvert[S any, D any](arr []S, f func(t S) D) []D {
	res := make([]D, len(arr))
	for i, t := range arr {
		res[i] = f(t)
	}
	return res
}

func SliceShuffle[T any](slice []T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func SliceRemove[T comparable](slice []T, e T) []T {
	for i, v := range slice {
		if v == e {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
