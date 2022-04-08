package mnemosyne

type index[K comparable, V any] struct{ kvstore map[K]*cacheEntry[V] }

func makeIndex[K comparable, V any]() *index[K, V] {
	return &index[K, V]{
		kvstore: make(map[K]*cacheEntry[V]),
	}
}

func (obj *index[K, V]) del(key K)                     { delete(obj.kvstore, key) }
func (obj *index[K, V]) set(key K, val *cacheEntry[V]) { obj.kvstore[key] = val }
func (obj *index[K, V]) get(key K) (val *cacheEntry[V], found bool) {
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

func (obj *index[K, V]) evict() {
	now := timeSource()
	var toDelete []K
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
