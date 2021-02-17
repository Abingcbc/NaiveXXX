package server

import (
	"NaiveDistributedMemCache/storage"
	"sync"
	"testing"
)

func TestModifyAfterCacheGet(t *testing.T) {
	cache := Cache{
		mu:      sync.Mutex{},
		storage: storage.NewLruCache(1000, nil),
	}
	cache.Put("test", []byte("test"))
	value := cache.Get("test")
	if string(value) != "test" {
		t.Error()
	}
	value = append(value, []byte("test")...)
	if string(cache.Get("test")) != "test" {
		t.Error()
	}
}
