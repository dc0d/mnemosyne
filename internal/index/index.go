package index

import (
	"github.com/dc0d/mnemosyne/internal/entry"
	"github.com/dc0d/mnemosyne/internal/support"
)

type Index[K comparable, V any] struct{ kvstore map[K]*entry.Entry[V] }

func NewIndex[K comparable, V any]() *Index[K, V] {
	return &Index[K, V]{
		kvstore: make(map[K]*entry.Entry[V]),
	}
}

func (obj *Index[K, V]) Del(key K)                      { delete(obj.kvstore, key) }
func (obj *Index[K, V]) Set(key K, val *entry.Entry[V]) { obj.kvstore[key] = val }
func (obj *Index[K, V]) Get(key K) (val *entry.Entry[V], found bool) {
	val, found = obj.kvstore[key]
	if !found {
		return
	}
	if val.Expired(support.TimeSource()) {
		obj.Del(key)
		return nil, false
	}
	return
}

func (obj *Index[K, V]) Evict() {
	now := support.TimeSource()
	var toDelete []K
	for key, val := range obj.kvstore {
		if !val.Expired(now) {
			continue
		}
		toDelete = append(toDelete, key)
	}
	for _, k := range toDelete {
		obj.Del(k)
	}
}
