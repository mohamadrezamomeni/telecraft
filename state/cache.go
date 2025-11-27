package state

import (
	"sync"
	"time"
)

type Cache struct {
	data         map[string]*State
	mutex        sync.RWMutex
	ttl          time.Duration
	defaultCache *Cache
	once         sync.Once
}

func New() *Cache {
	return &Cache{
		data: make(map[string]*State),
		ttl:  10 * time.Minute,
	}
}

func (c *Cache) Set(key string, state *State) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = state

	return nil
}

func (c *Cache) Get(key string) (*State, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	it, found := c.data[key]

	if !found || time.Now().After(it.Expiration) {
		return nil, false
	}
	return it, true
}

func (c *Cache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	return nil
}
