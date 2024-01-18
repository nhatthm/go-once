package once_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/once"
)

func TestFunc(t *testing.T) {
	t.Parallel()

	var actual atomic.Int64

	f := once.Func(func() {
		actual.Add(1)
	})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			f()
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(1), actual.Load())
}

func TestValue(t *testing.T) {
	t.Parallel()

	var (
		called atomic.Int64
		actual atomic.Int64
	)

	f := once.Value(func() int64 {
		called.Add(1)

		return called.Load()
	})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			actual.Add(f())
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(1), called.Load())
	assert.Equal(t, int64(100), actual.Load())
}

func TestValues(t *testing.T) {
	t.Parallel()

	var (
		called  atomic.Int64
		actual1 atomic.Int64
		actual2 atomic.Int64
	)

	f := once.Values(func() (int64, int64) {
		called.Add(1)

		return called.Load(), -called.Load()
	})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			v1, v2 := f()

			actual1.Add(v1)
			actual2.Add(v2)
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(1), called.Load())
	assert.Equal(t, int64(100), actual1.Load())
	assert.Equal(t, int64(-100), actual2.Load())
}

func TestFuncMap(t *testing.T) {
	t.Parallel()

	var actual atomic.Int64

	m := once.FuncMap[string]{}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			m.Do("key1", func() {
				actual.Add(1)
			})
		}()
	}

	wg.Wait()

	m.Do("key2", func() {
		actual.Add(1)
	})

	assert.Equal(t, int64(2), actual.Load())
	assert.Equal(t, 2, m.Len())

	m.Delete("key1")

	assert.Equal(t, 1, m.Len())

	m.Do("key1", func() {
		actual.Add(2)
	})

	assert.Equal(t, int64(4), actual.Load())
}

func TestFuncMap_Panics(t *testing.T) {
	t.Parallel()

	count := atomic.Int64{}
	m := once.FuncMap[string]{}
	f := func() {
		count.Add(1)
		panic("test")
	}

	for i := 0; i < 100; i++ {
		assert.Panics(t, func() {
			m.Do("key1", f)
		})
	}

	assert.Equal(t, int64(1), count.Load())
}

func TestValueMap(t *testing.T) {
	t.Parallel()

	var (
		actual atomic.Int64
		sum    atomic.Int64
	)

	m := once.ValueMap[string, int64]{}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sum.Add(m.Do("key1", func() int64 {
				actual.Add(1)

				return actual.Load()
			}))
		}()
	}

	wg.Wait()

	sum.Add(m.Do("key2", func() int64 {
		actual.Add(1)

		return 1
	}))

	assert.Equal(t, int64(2), actual.Load())
	assert.Equal(t, int64(101), sum.Load())
	assert.Equal(t, 2, m.Len())

	m.Delete("key1")

	assert.Equal(t, 1, m.Len())

	m.Do("key1", func() int64 {
		actual.Add(2)

		return actual.Load()
	})

	assert.Equal(t, int64(4), actual.Load())
}

func TestValueMap_Panics(t *testing.T) {
	t.Parallel()

	count := atomic.Int64{}
	m := once.ValueMap[string, string]{}
	f := func() string {
		count.Add(1)
		panic("test")
	}

	for i := 0; i < 100; i++ {
		assert.Panics(t, func() {
			m.Do("key1", f)
		})
	}

	assert.Equal(t, int64(1), count.Load())
}

func TestValueMap_Values(t *testing.T) {
	t.Parallel()

	m := once.ValueMap[string, string]{}

	m.Do("key1", func() string { return "value1" })
	m.Do("key2", func() string { return "value2" })

	actual := m.Values()
	expected := []string{"value1", "value2"}

	assert.ElementsMatch(t, expected, actual)
}

func TestValuesMap(t *testing.T) {
	t.Parallel()

	var (
		actual atomic.Int64
		sum1   atomic.Int64
		sum2   atomic.Int64
	)

	m := once.ValuesMap[string, int64, int64]{}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			v1, v2 := m.Do("key1", func() (int64, int64) {
				actual.Add(1)

				return actual.Load(), -actual.Load()
			})

			sum1.Add(v1)
			sum2.Add(v2)
		}()
	}

	wg.Wait()

	v1, v2 := m.Do("key2", func() (int64, int64) {
		actual.Add(1)

		return 1, -1
	})

	sum1.Add(v1)
	sum2.Add(v2)

	assert.Equal(t, int64(2), actual.Load())
	assert.Equal(t, int64(101), sum1.Load())
	assert.Equal(t, int64(-101), sum2.Load())
	assert.Equal(t, 2, m.Len())

	m.Delete("key1")

	assert.Equal(t, 1, m.Len())

	m.Do("key1", func() (int64, int64) {
		actual.Add(2)

		return 1, -1
	})

	assert.Equal(t, int64(4), actual.Load())
}

func TestValuesMap_Panics(t *testing.T) {
	t.Parallel()

	count := atomic.Int64{}
	m := once.ValuesMap[string, string, string]{}
	f := func() (string, string) {
		count.Add(1)
		panic("test")
	}

	for i := 0; i < 100; i++ {
		assert.Panics(t, func() {
			m.Do("key1", f)
		})
	}

	assert.Equal(t, int64(1), count.Load())
}

func TestPool_Get(t *testing.T) {
	t.Parallel()

	actual := atomic.Int64{}

	p := once.LazyValueMap[string, string]{
		New: func(key string) string {
			actual.Add(1)

			return key
		},
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			assert.Equal(t, "key1", p.Get("key1"))
		}()
	}

	wg.Wait()

	assert.Equal(t, "key2", p.Get("key2"))
	assert.Equal(t, int64(2), actual.Load())
	assert.Equal(t, 2, p.Len())

	p.Delete("key1")

	assert.Equal(t, 1, p.Len())
}
