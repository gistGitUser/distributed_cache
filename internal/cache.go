package internal

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value

	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val, ok := c.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", string(key))
	}
	return val, nil
}

func (c *Cache) Has(key []byte) (bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Delete(key []byte) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	delete(c.data, string(key))
	return nil
}