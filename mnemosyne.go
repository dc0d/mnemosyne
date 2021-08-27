package mnemosyne

import (
	"sync"
	"time"
)

type Cache struct {
	index *index
	lock  sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		index: makeIndex(),
	}
}

func (obj *Cache) GetOrInit(key string, init Initializer) TV {
	obj.lock.Lock()

	entry, found := obj.index.get(key)
	if !found {
		entry = &cacheEntry{}
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

func (obj *Cache) Put(key string, init Initializer) TV {
	obj.lock.Lock()

	entry := &cacheEntry{}
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

func (obj *Cache) Remove(firstKey string, rest ...string) {
	rest = append(rest, firstKey)

	obj.lock.Lock()
	defer obj.lock.Unlock()

	for _, key := range rest {
		obj.index.del(key)
	}
}

func (obj *Cache) Get(key string) (val TV, found bool) {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	entry, found := obj.index.get(key)
	if !found {
		return
	}

	return entry.value, true
}

func (obj *Cache) Evict() {
	obj.lock.Lock()
	defer obj.lock.Unlock()

	obj.index.evict()
}

type Initializer func() (TV, time.Duration)

//

type index struct{ kvstore map[TK]*cacheEntry }

func makeIndex() *index {
	return &index{
		kvstore: make(map[TK]*cacheEntry),
	}
}

func (obj *index) del(key TK)                  { delete(obj.kvstore, key) }
func (obj *index) set(key TK, val *cacheEntry) { obj.kvstore[key] = val }
func (obj *index) get(key TK) (val *cacheEntry, found bool) {
	val, found = obj.kvstore[key]
	if !found {
		return
	}
	if val.expired(timeSource()) {
		obj.del(key)
		return nil, false
	}
	return
}

func (obj *index) evict() {
	now := timeSource()
	var toDelete []TK
	for key, val := range obj.kvstore {
		if !val.expired(now) {
			continue
		}
		toDelete = append(toDelete, key)
	}
	for _, k := range toDelete {
		obj.del(k)
	}
}

//

type cacheEntry struct {
	value    TV
	once     sync.Once
	deadline time.Time
}

func (obj *cacheEntry) expires() bool              { return !obj.deadline.IsZero() }
func (obj *cacheEntry) expired(now time.Time) bool { return obj.expires() && obj.deadline.Before(now) }

type (
	TK = string
	TV = struct{ Message string }
)

var (
	timeSource = time.Now
)
