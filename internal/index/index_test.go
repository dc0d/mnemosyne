package index

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dc0d/mnemosyne/internal/entry"
	"github.com/dc0d/mnemosyne/internal/support"
	"github.com/dc0d/mnemosyne/internal/testsupport"
)

func Test_index(t *testing.T) {
	t.Run(`should evict entry on get`, PreserveTimeSource(func(t *testing.T) {
		support.TimeSource = time.Now

		sut := NewIndex[string, Value]()

		e := entry.New[Value]()
		te := entry.NewTestEntry[Value](e)
		te.SetDeadline(time.Now().Add(time.Hour * -1))

		sut.Set("1", e)

		_, found := sut.Get("1")

		assert.False(t, found)
	}))

	t.Run(`should evict expired antries when calling evict`, PreserveTimeSource(func(t *testing.T) {
		support.TimeSource = time.Now

		sut := NewIndex[string, Value]()

		// create 10 expired entries
		for i := 0; i < 10; i++ {
			e := entry.New[Value]()
			te := entry.NewTestEntry[Value](e)
			te.SetDeadline(time.Now().Add(time.Hour * -1))

			sut.Set(fmt.Sprintf("expired-%v", i), e)
		}

		// create entries that expire in the future
		for i := 0; i < 10; i++ {
			e := entry.New[Value]()
			te := entry.NewTestEntry[Value](e)
			te.SetDeadline(time.Now().Add(time.Hour * 24))

			sut.Set(fmt.Sprintf("future-%v", i), e)
		}

		sut.Evict()

		for i := 0; i < 10; i++ {
			_, found := sut.Get(fmt.Sprintf("expired-%v", i))
			assert.False(t, found)
		}

		for i := 0; i < 10; i++ {
			_, found := sut.Get(fmt.Sprintf("future-%v", i))
			assert.True(t, found)
		}
	}))
}

type Value = testsupport.Value

var PreserveTimeSource = testsupport.PreserveTimeSource
