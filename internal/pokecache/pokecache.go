package pokecache

import (
	"sync"
	"time"
)

func New(interval time.Duration) Cache {
	entries := make(map[string]cacheEntry)
	mux := &sync.Mutex{}
	cache := Cache{entries, mux}

	go cache.reapLoop(interval)

	return cache
}

type Cache struct {
	entries map[string]cacheEntry
	// Could use an RWMutex, but reads basically only happen on user input
	mux *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	createdAt := time.Now()
	entry := cacheEntry{
		createdAt,
		val,
	}
	c.mux.Lock()
	c.entries[key] = entry
	c.mux.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	entry, ok := c.entries[key]
	c.mux.Unlock()

	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for time := range ticker.C {
		c.cleanBefore(time.Add(-interval))
	}
}

func (c *Cache) cleanBefore(time time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for key, entry := range c.entries {
		if entry.createdAt.Before(time) {
			delete(c.entries, key)
		}
	}
}
