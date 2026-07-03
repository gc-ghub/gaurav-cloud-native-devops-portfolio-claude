package cache

import (
	"sync"
	"time"
)

// Cache holds one TTL-bound value. GetOrFetch serves the cached value while
// fresh; on a stale/empty cache it calls fetch(). If fetch fails but a prior
// value exists (even stale), that value is returned instead of the error —
// callers only see an error on a cold cache with a failing upstream.
type Cache[T any] struct {
	mu        sync.RWMutex
	value     T
	hasValue  bool
	fetchedAt time.Time
	ttl       time.Duration
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{ttl: ttl}
}

func (c *Cache[T]) GetOrFetch(fetch func() (T, error)) (T, error) {
	c.mu.RLock()
	fresh := c.hasValue && time.Since(c.fetchedAt) < c.ttl
	value := c.value
	c.mu.RUnlock()

	if fresh {
		return value, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Another goroutine may have refreshed while we waited for the lock.
	if c.hasValue && time.Since(c.fetchedAt) < c.ttl {
		return c.value, nil
	}

	newValue, err := fetch()
	if err != nil {
		if c.hasValue {
			return c.value, nil
		}
		var zero T
		return zero, err
	}

	c.value = newValue
	c.hasValue = true
	c.fetchedAt = time.Now()
	return c.value, nil
}
