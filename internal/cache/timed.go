// Package cache contains a single implementation of a cache with a TTL. We can
// keep it mega simple because only 40~ users could even be put in here, but I've
// got a basic cleanup loop set up anyway.
package cache

import (
	"context"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]cacheItem[V]
	ttl   time.Duration
}

type cacheItem[V any] struct {
	value     V
	expiresAt time.Time
}

func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]cacheItem[V]),
		ttl:   ttl,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		var zero V
		return zero, false
	}
	return item.value, true
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem[V]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache[K, V]) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.mu.Lock()
				for k, item := range c.items {
					if time.Now().After(item.expiresAt) {
						delete(c.items, k)
					}
				}
				c.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}
