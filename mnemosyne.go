package mnemosyne

import (
	"sync"

	"github.com/dc0d/mnemosyne/internal/entry"
	"github.com/dc0d/mnemosyne/internal/index"
)

type Cache[K comparable, V any] struct {
	index *index.Index[K, V]
	lock  sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		index: index.NewIndex[K, V](),
	}
}

func (obj *Cache[K, V]) GetOrInit(key K, init Initializer[V]) V {
	obj.lock.Lock()

	e, found := obj.index.Get(key)
	if !found {
		e = entry.New[V]()
		obj.index.Set(key, e)
	}

	obj.lock.Unlock()

	e.Apply(init)

	return e.Value()
}

func (obj *Cache[K, V]) Put(key K, init Initializer[V]) V {
	obj.lock.Lock()

	e := entry.New[V]()
	obj.index.Set(key, e)

	obj.lock.Unlock()

	e.Apply(init)

	return e.Value()
}

func (obj *Cache[K, V]) Remove(firstKey K, rest ...K) {
	rest = append(rest, firstKey)

	obj.lock.Lock()
	defer obj.lock.Unlock()

	for _, key := range rest {
		obj.index.Del(key)
	}
}

func (obj *Cache[K, V]) Get(key K) (val V, found bool) {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	e, found := obj.index.Get(key)
	if !found {
		return
	}

	return e.Value(), true
}

func (obj *Cache[K, V]) Evict() {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	obj.index.Evict()
}
