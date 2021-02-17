package server

import (
	"NaiveDistributedMemCache/storage"
	"sync"
)

type Cache struct {
	mu sync.Mutex
	storage storage.Storage
}

func NewCache() *Cache {
	return &Cache{
		mu:      sync.Mutex{},
		storage: storage.NewLruCache(2<<10, nil),
	}
}

func (cache *Cache) Put(key string, bytes []byte) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	return cache.storage.Put(key, storage.Value{Bytes: bytes})
}

func (cache *Cache) Get(key string) []byte {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if value, ok := cache.storage.Get(key); ok {
		return value.Bytes
	}
	return nil
}
