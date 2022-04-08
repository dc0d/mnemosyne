package mnemosyne

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	index *index[K, V]
	lock  sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		index: makeIndex[K, V](),
	}
}

func (obj *Cache[K, V]) GetOrInit(key K, init Initializer[V]) V {
	obj.lock.Lock()

	entry, found := obj.index.get(key)
	if !found {
		entry = &cacheEntry[V]{}
		obj.index.set(key, entry)
	}

	obj.lock.Unlock()

	entry.once.Do(func() {
		v, d := init()
		entry.value = v
		if d <= 0 {
			return
		}
		entry.deadline = timeSource().Add(d)
	})

	return entry.value
}

func (obj *Cache[K, V]) Put(key K, init Initializer[V]) V {
	obj.lock.Lock()

	entry := &cacheEntry[V]{}
	obj.index.set(key, entry)

	obj.lock.Unlock()

	entry.once.Do(func() {
		v, d := init()
		entry.value = v
		if d <= 0 {
			return
		}
		entry.deadline = timeSource().Add(d)
	})

	return entry.value
}

func (obj *Cache[K, V]) Remove(firstKey K, rest ...K) {
	rest = append(rest, firstKey)

	obj.lock.Lock()
	defer obj.lock.Unlock()

	for _, key := range rest {
		obj.index.del(key)
	}
}

func (obj *Cache[K, V]) Get(key K) (val V, found bool) {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	entry, found := obj.index.get(key)
	if !found {
		return
	}

	return entry.value, true
}

func (obj *Cache[K, V]) Evict() {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	obj.index.evict()
}
