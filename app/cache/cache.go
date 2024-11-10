package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	data       string
	expiration *time.Time
}

func (ci *cacheItem) IsExpired() bool {
	if ci.expiration == nil {
		return false
	}

	return time.Now().After(*ci.expiration)
}

type Cache struct {
	cache map[string]cacheItem
	m     sync.Mutex
}

func CreateCache() *Cache {
	cache := &Cache{
		cache: make(map[string]cacheItem),
		m:     sync.Mutex{},
	}

	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	val, isExist := c.cache[key]

	if !isExist {
		return "", false
	}

	if val.IsExpired() {
		delete(c.cache, key)
		return "", false
	}

	return val.data, isExist
}

func (c *Cache) Set(key string, value string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = cacheItem{
		data:       value,
		expiration: nil,
	}
}

// Sets cache with expiration time in milliseconds
func (c *Cache) SetWithExpiration(key string, value string, expiration int) {
	c.m.Lock()
	defer c.m.Unlock()

	expirationTime := time.Now().Add(time.Duration(expiration) * time.Millisecond)

	c.cache[key] = cacheItem{
		data:       value,
		expiration: &expirationTime,
	}
}
