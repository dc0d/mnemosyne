package mnemosyne

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dc0d/mnemosyne/internal/support"
	"github.com/dc0d/mnemosyne/internal/testsupport"
)

//nolint:funlen
func TestCache(t *testing.T) {
	t.Run(`should cache when calling GetOrInit`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		expectedValue := val1
		actualValue := sut.GetOrInit(key1, func() (Value, time.Duration) {
			return expectedValue, 0
		})

		assert.Equal(t, expectedValue, actualValue)
	}))

	t.Run(`should not overwrite after cached by calling GetOrInit`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		initialValue := val1
		actualValue := sut.GetOrInit(key1, func() (Value, time.Duration) {
			return initialValue, 0
		})

		assert.Equal(t, initialValue, actualValue)

		changedValue := val2
		actualValue = sut.GetOrInit(key1, func() (Value, time.Duration) {
			return changedValue, 0
		})

		assert.Equal(t, initialValue, actualValue)
	}))

	t.Run(`should evict expired entry by calling Get - initialized by GetOrInit`, PreserveTimeSource(func(t *testing.T) {
		support.TimeSource = time.Now
		sut := New[string, Value]()

		initialValue := val1
		sut.GetOrInit(key1, func() (Value, time.Duration) {
			return initialValue, time.Millisecond
		})

		assert.Eventually(t, func() bool {
			result, _ := sut.Get(key1)
			return result == Value{}
		}, time.Millisecond*300, time.Millisecond*20)

		assert.Eventually(t, func() bool {
			_, found := sut.Get(key1)
			return !found
		}, time.Millisecond*300, time.Millisecond*20)
	}))

	t.Run(`should cache when calling Put`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		expectedValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return expectedValue, 0
		})
		actual, _ := sut.index.Get(key1)

		assert.Equal(t, expectedValue, actual.Value())
	}))

	t.Run(`should overwrite after cached by calling Put`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		initialValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return initialValue, 0
		})
		actual, _ := sut.index.Get(key1)

		assert.Equal(t, initialValue, actual.Value())

		changedValue := val2
		sut.Put(key1, func() (Value, time.Duration) {
			return changedValue, 0
		})
		actual, _ = sut.index.Get(key1)

		assert.Equal(t, changedValue, actual.Value())
	}))

	t.Run(`should evict expired entry by calling Get - initialized by Put`, PreserveTimeSource(func(t *testing.T) {
		support.TimeSource = time.Now
		sut := New[string, Value]()

		initialValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return initialValue, time.Millisecond
		})

		assert.Eventually(t, func() bool {
			_, found := sut.Get(key1)
			return !found
		}, time.Millisecond*300, time.Millisecond*20)

		assert.Eventually(t, func() bool {
			result, _ := sut.Get(key1)
			return result == Value{}
		}, time.Millisecond*300, time.Millisecond*20)
	}))

	t.Run(`should return not found when not cached when calling Get`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		_, found := sut.Get(key1)

		assert.False(t, found)
	}))

	t.Run(`should find the entry when already cached when calling Get`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		expectedValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)
	}))

	t.Run(`should return not found when entry is removed when calling Get`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		expectedValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)

		sut.Remove(key1)

		val, found = sut.Get(key1)

		assert.False(t, found)
		assert.Equal(t, Value{}, val)
	}))

	t.Run(`should be fine calling Remove with a non-existing key`, PreserveTimeSource(func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		support.TimeSource = func() time.Time { return now }
		sut := New[string, Value]()

		sut.Remove(key1)

		expectedValue := val1
		sut.Put(key1, func() (Value, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)

		sut.Remove(key1)

		val, found = sut.Get(key1)

		assert.False(t, found)
		assert.Equal(t, Value{}, val)
	}))

	t.Run(`should evict expired entries`, PreserveTimeSource(func(t *testing.T) {
		support.TimeSource = time.Now
		sut := New[string, Value]()

		for i := 1; i <= 10; i++ {
			key := fmt.Sprint(i)
			sut.Put(key, func() (Value, time.Duration) {
				return val1, time.Nanosecond
			})
		}

		for i := 100; i <= 110; i++ {
			key := fmt.Sprint(i)
			sut.Put(key, func() (Value, time.Duration) {
				return val1, time.Hour * 24 * 7
			})
		}

		assert.Eventually(t, func() bool {
			sut.Evict()
			result := true
			for i := 1; i <= 10; i++ {
				key := fmt.Sprint(i)
				_, found := sut.Get(key)
				result = result && !found
			}
			return result
		}, time.Millisecond*300, time.Millisecond*20)
	}))
}

const (
	key1 = "a-key"
)

var (
	val1 = Value{Message: "MSG 1"}
	val2 = Value{Message: "MSG 2"}
)

type Value = testsupport.Value

var PreserveTimeSource = testsupport.PreserveTimeSource
