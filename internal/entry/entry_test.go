package entry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dc0d/mnemosyne/internal/testsupport"
)

func Test_Entry(t *testing.T) {
	t.Run(`expires`, func(t *testing.T) {
		sut := New[testsupport.Value]()

		assert.False(t, sut.expires())
	})

	t.Run(`Expired`, func(t *testing.T) {
		sut := New[testsupport.Value]()

		assert.False(t, sut.Expired(time.Now()))
	})

	t.Run(`Apply - no TTL`, func(t *testing.T) {
		sut := New[testsupport.Value]()

		sut.Apply(func() (testsupport.Value, time.Duration) {
			return testsupport.Value{Message: "hello"}, 0
		})

		sut.Apply(func() (testsupport.Value, time.Duration) {
			return testsupport.Value{Message: "we should not see this value"}, 0
		})

		assert.Equal(t, "hello", sut.Value().Message)
	})

	t.Run(`Apply - with TTL`, func(t *testing.T) {
		sut := New[testsupport.Value]()

		sut.Apply(func() (testsupport.Value, time.Duration) {
			return testsupport.Value{Message: ""}, time.Second * 5
		})

		assert.InEpsilon(t,
			time.Now().Add(time.Second*5).UnixNano(),
			sut.deadline.UnixNano(),
			float64(time.Millisecond*500))
	})
}
