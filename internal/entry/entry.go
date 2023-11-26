package entry

import (
	"sync"
	"time"

	"github.com/dc0d/mnemosyne/internal/support"
)

type Entry[V any] struct {
	value    V
	once     sync.Once
	deadline time.Time
}

func New[V any]() *Entry[V] {
	return &Entry[V]{}
}

func (obj *Entry[V]) Value() V {
	return obj.value
}

func (obj *Entry[V]) Expired(t time.Time) bool {
	return obj.expires() && obj.deadline.Before(t)
}

func (obj *Entry[V]) expires() bool { return !obj.deadline.IsZero() }

func (obj *Entry[V]) Apply(init func() (V, time.Duration)) {
	obj.once.Do(func() {
		v, d := init()
		obj.value = v
		if d <= 0 {
			return
		}
		obj.deadline = support.TimeSource().Add(d)
	})
}

type TestEntry[V any] struct {
	*Entry[V]
}

func NewTestEntry[V any](e *Entry[V]) *TestEntry[V] {
	return &TestEntry[V]{Entry: e}
}

func (obj *TestEntry[V]) SetDeadline(t time.Time) {
	obj.Entry.deadline = t
}
