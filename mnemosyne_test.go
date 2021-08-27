package mnemosyne

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// //nolint:funlen
func TestCache(t *testing.T) {
	t.Run(`should cache when calling GetOrInit`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		expectedValue := val1
		actualValue := sut.GetOrInit(key1, func() (TV, time.Duration) {
			return expectedValue, 0
		})

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run(`should not overwrite after cached by calling GetOrInit`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		initialValue := val1
		actualValue := sut.GetOrInit(key1, func() (TV, time.Duration) {
			return initialValue, 0
		})

		assert.Equal(t, initialValue, actualValue)

		changedValue := val2
		actualValue = sut.GetOrInit(key1, func() (TV, time.Duration) {
			return changedValue, 0
		})

		assert.Equal(t, initialValue, actualValue)
	})

	t.Run(`should evict expired entry by calling Get - initialized by GetOrInit`, func(t *testing.T) {
		timeSource = time.Now
		sut := NewCache()

		initialValue := val1
		sut.GetOrInit(key1, func() (TV, time.Duration) {
			return initialValue, time.Millisecond
		})

		assert.Eventually(t, func() bool {
			result, _ := sut.Get(key1)
			return result == zeroData
		}, time.Millisecond*300, time.Millisecond*20)

		assert.Eventually(t, func() bool {
			_, found := sut.Get(key1)
			return !found
		}, time.Millisecond*300, time.Millisecond*20)
	})

	t.Run(`should cache when calling Put`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		expectedValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return expectedValue, 0
		})

		assert.Equal(t, expectedValue, sut.index.kvstore[key1].value)
	})

	t.Run(`should overwrite after cached by calling Put`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		initialValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return initialValue, 0
		})

		assert.Equal(t, initialValue, sut.index.kvstore[key1].value)

		changedValue := val2
		sut.Put(key1, func() (TV, time.Duration) {
			return changedValue, 0
		})

		assert.Equal(t, changedValue, sut.index.kvstore[key1].value)
	})

	t.Run(`should evict expired entry by calling Get - initialized by Put`, func(t *testing.T) {
		timeSource = time.Now
		sut := NewCache()

		initialValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return initialValue, time.Millisecond
		})

		assert.Eventually(t, func() bool {
			_, found := sut.Get(key1)
			return !found
		}, time.Millisecond*300, time.Millisecond*20)

		assert.Eventually(t, func() bool {
			result, _ := sut.Get(key1)
			return result == zeroData
		}, time.Millisecond*300, time.Millisecond*20)
	})

	t.Run(`should return not found when not cached when calling Get`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		_, found := sut.Get(key1)

		assert.False(t, found)
	})

	t.Run(`should find the entry when already cached when calling Get`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		expectedValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)
	})

	t.Run(`should return not found when entry is removed when calling Get`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		expectedValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)

		sut.Remove(key1)

		val, found = sut.Get(key1)

		assert.False(t, found)
		assert.Equal(t, zeroData, val)
	})

	t.Run(`should be fine calling Remove with a non-existing key`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		sut := NewCache()

		sut.Remove(key1)

		expectedValue := val1
		sut.Put(key1, func() (TV, time.Duration) {
			return expectedValue, 0
		})

		val, found := sut.Get(key1)

		assert.True(t, found)
		assert.Equal(t, expectedValue, val)

		sut.Remove(key1)

		val, found = sut.Get(key1)

		assert.False(t, found)
		assert.Equal(t, zeroData, val)
	})

	t.Run(`should evict expired entries`, func(t *testing.T) {
		timeSource = time.Now
		sut := NewCache()

		for i := 1; i <= 10; i++ {
			key := fmt.Sprint(i)
			sut.Put(key, func() (TV, time.Duration) {
				return val1, time.Nanosecond
			})
		}

		for i := 100; i <= 110; i++ {
			key := fmt.Sprint(i)
			sut.Put(key, func() (TV, time.Duration) {
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
	})
}

const (
	key1 = "a-key"
)

var (
	zeroData struct{ Message string }
	val1     = struct{ Message string }{Message: "MSG 1"}
	val2     = struct{ Message string }{Message: "MSG 2"}
)

func Test_index(t *testing.T) {
	t.Run(`should evict on get`, func(t *testing.T) {
		now := time.Date(2021, 8, 25, 20, 30, 0, 0, time.Local)
		timeSource = func() time.Time { return now }
		const minutes = 10
		sut := makeIndex()
		expected := makeIndex()
		for i := -1 * minutes; i <= minutes; i++ {
			deadline := now.Add(time.Minute * time.Duration(i))
			key := fmt.Sprint(i)
			entry := &cacheEntry{deadline: deadline}
			sut.set(key, entry)

			if deadline.Before(now) {
				continue
			}

			expected.set(key, entry)
		}

		for i := -1 * minutes; i <= minutes; i++ {
			key := fmt.Sprint(i)

			sut.get(key)
		}

		want := fmt.Sprint(expected)
		got := fmt.Sprint(sut)

		assert.Equal(t, want, got)
	})
}

func Test_cacheEntry(t *testing.T) {
	t.Run(`expires`, func(t *testing.T) {
		sut := &cacheEntry{}

		assert.False(t, sut.expires())
	})

	t.Run(`expired`, func(t *testing.T) {
		sut := &cacheEntry{}

		aTime := time.Now()

		assert.False(t, sut.expired(aTime))
	})
}
