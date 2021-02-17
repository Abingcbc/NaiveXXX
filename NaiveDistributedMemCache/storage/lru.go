package storage

import "container/list"

type LruCache struct {
	maxBytes int64
	currentBytes int64
	linkedList *list.List
	indexMap map[string]*list.Element
	onEvicted func(key string, value Value)
}

func NewLruCache(maxBytes int64, onEvicted func(key string, value Value)) *LruCache {
	return &LruCache{
		maxBytes:     maxBytes,
		currentBytes: 0,
		linkedList:   list.New(),
		indexMap:     make(map[string]*list.Element),
		onEvicted:    onEvicted,
	}
}

func (lru *LruCache) Put(key string, value Value) bool {
	if value.Len() > lru.maxBytes {
		return false
	}
	if element, ok := lru.indexMap[key]; ok {
		entry := element.Value.(*Entry)
		lru.currentBytes += value.Len() - entry.value.Len()
		lru.linkedList.MoveToFront(element)
		entry.value = value
	} else {
		lru.currentBytes += int64(len(key)) + value.Len()
		lru.linkedList.PushFront(&Entry{
			key:   key,
			value: value,
		})
		lru.indexMap[key] = lru.linkedList.Front()
	}

	for {
		if lru.currentBytes <= lru.maxBytes {
			break
		}
		staleCache := lru.linkedList.Back()
		lru.linkedList.Remove(staleCache)
		entry := staleCache.Value.(*Entry)
		delete(lru.indexMap, entry.key)
		if lru.onEvicted != nil {
			lru.onEvicted(entry.key, entry.value)
		}

		lru.currentBytes -= int64(len(entry.key)) + entry.value.Len()
	}
	return true
}

func (lru *LruCache) Get(key string) (Value, bool) {
	if element, ok := lru.indexMap[key]; ok {
		lru.linkedList.MoveToFront(element)
		return element.Value.(*Entry).value, true
	}
	return Value{}, false
}
