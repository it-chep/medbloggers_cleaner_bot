package inmemory_cache

import (
	"sync"
	"time"
)

type cacheValue[Val any] struct {
	value    Val
	ttl      time.Duration
	createAt time.Time
}

type Cache[Key comparable, Val any] struct {
	cache map[Key]*cacheValue[Val]
	mutex sync.RWMutex
}

func NewCache[Key comparable, Val any](cleanPeriod time.Duration) *Cache[Key, Val] {
	cache := &Cache[Key, Val]{
		cache: make(map[Key]*cacheValue[Val]),
	}

	go func(cache *Cache[Key, Val]) {
		timer := time.NewTicker(cleanPeriod)

		for tick := range timer.C {
			cache.cleaner(tick.UTC())
		}

	}(cache)

	return cache
}

func (cache *Cache[Key, Val]) Set(key Key, value Val, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.cache[key] = &cacheValue[Val]{
		value:    value,
		createAt: time.Now().UTC(),
		ttl:      ttl,
	}
}

func (cache *Cache[Key, Val]) Get(key Key) (Val, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	val, ok := cache.cache[key]

	if !ok {
		var zeroValue Val
		return zeroValue, false
	}

	return val.value, ok
}

func (cache *Cache[Key, Val]) Del(key Key) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.cache, key)
}

func (cache *Cache[Key, Val]) cleaner(now time.Time) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	removeKeys := make([]Key, 0)

	for key, val := range cache.cache {
		if now.After(val.createAt.Add(val.ttl)) {
			removeKeys = append(removeKeys, key)
		}
	}

	for _, key := range removeKeys {
		delete(cache.cache, key)
	}
}
