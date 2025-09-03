package cache

import (
	"log"
	"sync"
	"time"
)

type InMemoryCache struct {
	ttl   time.Duration
	store sync.Map
}

type cacheData struct {
	createdAt time.Time
	value     interface{}
}

func NewInMemory(ttl, clearTicker time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		ttl:   ttl,
		store: sync.Map{},
	}

	go func(cache *InMemoryCache) {
		ticker := time.NewTicker(clearTicker)
		for {
			<-ticker.C
			log.Println("running cache cleaning worker")
			if cache == nil {
				return
			}
			cache.store.Range(func(k, v interface{}) bool {
				cd, ok := v.(cacheData)
				if !ok {
					log.Println("wrong cache data type detected in cache store")
					cache.store.Delete(k)
					return true
				}
				if time.Since(cd.createdAt) > cache.ttl {
					cache.store.Delete(k)
				}
				return true
			})
		}
	}(cache)

	return cache
}

func (imc *InMemoryCache) Set(key string, value interface{}) error {
	cd := cacheData{
		createdAt: time.Now(),
		value:     value,
	}
	imc.store.Store(key, cd)
	return nil
}

func (imc *InMemoryCache) Get(key string) (interface{}, error) {
	v, ok := imc.store.Load(key)
	if !ok {
		return nil, ErrCacheMiss
	}

	cd, ok := v.(cacheData)
	if !ok {
		return nil, ErrInvalidCacheValue
	}

	if time.Since(cd.createdAt) > imc.ttl {
		imc.store.Delete(key)
		return nil, ErrTTLExpired
	}

	return cd.value, nil
}
