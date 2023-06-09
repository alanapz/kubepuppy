package utils

import (
	"github.com/google/uuid"
	"reflect"
	"strings"
	"sync"
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

func MapWithError[T any, U any](src []T, mapper func(T) (U, error)) ([]U, error) {
	dest := make([]U, 0)

	for _, v := range src {
		r, err := mapper(v)

		if err != nil {
			return nil, err
		}

		dest = append(dest, r)
	}

	return dest, nil
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

func MapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k, _ := range m {
		r = append(r, k)
	}
	return r
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

const mutexLocked = 1

func IsMutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func IsRWMutexWriteLocked(rw *sync.RWMutex) bool {
	// RWMutex has a "w" sync.Mutex field for write lock
	state := reflect.ValueOf(rw).Elem().FieldByName("w").FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func IsRWMutexReadLocked(rw *sync.RWMutex) bool {
	return reflect.ValueOf(rw).Elem().FieldByName("readerCount").Int() > 0
}
