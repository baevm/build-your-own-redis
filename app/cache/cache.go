package cache

import "sync"

type Cache struct {
	cache map[string]string
	m     sync.Mutex
}

func CreateCache() *Cache {
	cache := &Cache{
		cache: make(map[string]string),
		m:     sync.Mutex{},
	}

	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	val, isExist := c.cache[key]
	return val, isExist
}

func (c *Cache) Set(key string, value string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = value
}
