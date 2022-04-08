package mnemosyne

import (
	"sync"
	"time"
)

type cacheEntry[V any] struct {
	value    V
	once     sync.Once
	deadline time.Time
}

func (obj *cacheEntry[V]) expires() bool { return !obj.deadline.IsZero() }
func (obj *cacheEntry[V]) expired(now time.Time) bool {
	return obj.expires() && obj.deadline.Before(now)
}
