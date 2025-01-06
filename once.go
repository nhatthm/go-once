// Package once provides a simple way to do something exactly once.
package once

import (
	"sync"
)

// Once is an object that will perform exactly one action.
//
// A Once must not be copied after first use.
//
// In the terminology of the Go memory model,
// the return from f “synchronizes before”
// the return from any call of once.Do(f).
type Once = sync.Once

// Func returns a function that invokes f only once. The returned function
// may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func Func(f func()) func() {
	return sync.OnceFunc(f)
}

// Value returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func Value[T any](f func() T) func() T {
	return sync.OnceValue(f)
}

// Values returns a function that invokes f only once and returns the values
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func Values[T1, T2 any](f func() (T1, T2)) func() (T1, T2) {
	return sync.OnceValues(f)
}

// Map is simpler version of sync.Map.
type Map[K, V any] sync.Map

// Get returns the value for key if it exists. Otherwise, it returns default value.
func (m *Map[K, V]) Get(key K, defaultValue V) V {
	v, _ := (*sync.Map)(m).LoadOrStore(key, defaultValue)

	return v.(V) //nolint: errcheck
}

// Delete removes the key and its value from the map.
func (m *Map[K, V]) Delete(key K) {
	(*sync.Map)(m).Delete(key)
}

// Len returns the number of entries in the map.
func (m *Map[K, V]) Len() int {
	result := 0

	(*sync.Map)(m).Range(func(_, _ any) bool {
		result++

		return true
	})

	return result
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	(*sync.Map)(m).Range(func(key, value any) bool {
		return f(key.(K), value.(V)) //nolint: errcheck
	})
}

// FuncMap is a map of functions that are only computed once.
type FuncMap[K any] Map[K, func()]

// Do executes the function only once for the given key.
//
// If f panics, the function will panic with the same value on every call.
func (m *FuncMap[K]) Do(key K, f func()) {
	(*Map[K, func()])(m).Get(key, sync.OnceFunc(f))()
}

// Delete deletes the key and its function from the map.
func (m *FuncMap[K]) Delete(key K) {
	(*Map[K, func()])(m).Delete(key)
}

// Len returns the number of functions in the map.
func (m *FuncMap[K]) Len() int {
	return (*Map[K, func()])(m).Len()
}

// ValueMap is a map of functions that are only computed once.
type ValueMap[K, V any] Map[K, func() V]

// Do executes the function only once for the given key and returns the value returned by f.
//
// If f panics, the function will panic with the same value on every call.
func (m *ValueMap[K, V]) Do(key K, f func() V) V {
	return (*Map[K, func() V])(m).Get(key, sync.OnceValue(f))()
}

// Delete removes the key and its function from the map.
func (m *ValueMap[K, V]) Delete(key K) {
	(*Map[K, func() V])(m).Delete(key)
}

// Len returns the number of functions in the map.
func (m *ValueMap[K, V]) Len() int {
	return (*Map[K, func() V])(m).Len()
}

// Values returns all values in the map. If f panics, the function will panic with the same value on every call.
func (m *ValueMap[K, V]) Values() []V {
	result := make([]V, 0)

	(*Map[K, func() V])(m).Range(func(_ K, f func() V) bool {
		result = append(result, f())

		return true
	})

	return result
}

// ValuesMap is a map of functions that are only computed once.
type ValuesMap[K, V1, V2 any] Map[K, func() (V1, V2)]

// Do executes the function only once for the given key and returns the values returned by f.
//
// If f panics, the function will panic with the same value on every call.
func (m *ValuesMap[K, V1, V2]) Do(key K, f func() (V1, V2)) (V1, V2) {
	return (*Map[K, func() (V1, V2)])(m).Get(key, sync.OnceValues(f))()
}

// Delete removes the key and its function from the map.
func (m *ValuesMap[K, V1, V2]) Delete(key K) {
	(*Map[K, func() (V1, V2)])(m).Delete(key)
}

// Len returns the number of functions in the map.
func (m *ValuesMap[K, V1, V2]) Len() int {
	return (*Map[K, func() (V1, V2)])(m).Len()
}

// LazyValueMap initializes values when they are first accessed.
type LazyValueMap[K, V any] struct {
	New func(key K) V

	entries ValueMap[K, V]
}

// Get returns the value for key if it exists. Otherwise, it calls New and stores.
func (p *LazyValueMap[K, V]) Get(key K) V {
	return p.entries.Do(key, func() V {
		return p.New(key)
	})
}

// Delete removes the key and its value from the map.
func (p *LazyValueMap[K, V]) Delete(key K) {
	p.entries.Delete(key)
}

// Len returns the number of entries in the map.
func (p *LazyValueMap[K, V]) Len() int {
	return p.entries.Len()
}
