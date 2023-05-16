package utils

import (
	"github.com/google/uuid"
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Guid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func Ptr[T any](x T) *T {
	return &x
}

func Map[T any, U any](src []T, mapper func(T) U) []U {
	dest := make([]U, len(src))
	for k, v := range src {
		dest[k] = mapper(v)
	}
	return dest
}

func Min[T ~int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T ~int](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func MapValues[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func Filter[T any](src []T, predicate func(T) bool) []T {
	r := make([]T, 0, 0)
	for _, v := range src {
		if predicate(v) {
			r = append(r, v)
		}
	}
	return r
}

func Any[T any](src []T, predicate func(T) bool) bool {
	for _, v := range src {
		if predicate(v) {
			return true
		}
	}
	return false
}
