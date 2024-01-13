package once_test

import (
	"testing"

	"go.nhat.io/once"
)

func BenchmarkFunc(b *testing.B) {
	f := once.Func(func() {})

	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkValue(b *testing.B) {
	f := once.Value(func() int {
		return 0
	})

	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkValues(b *testing.B) {
	f := once.Values(func() (int, int) {
		return 0, 0
	})

	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkFuncMap_Do(b *testing.B) {
	m := once.FuncMap[int]{}

	for i := 0; i < b.N; i++ {
		m.Do(i, func() {})
	}
}

func BenchmarkValueMap_Do(b *testing.B) {
	m := once.ValueMap[int, int]{}

	for i := 0; i < b.N; i++ {
		m.Do(i, func() int {
			return 0
		})
	}
}

func BenchmarkValuesMap_Do(b *testing.B) {
	m := once.ValuesMap[int, int, int]{}

	for i := 0; i < b.N; i++ {
		m.Do(i, func() (int, int) {
			return 0, 0
		})
	}
}

func BenchmarkPool_Get(b *testing.B) {
	p := once.LazyValueMap[int, int]{
		New: func(key int) int {
			return key
		},
	}

	for i := 0; i < b.N; i++ {
		p.Get(i)
	}
}
