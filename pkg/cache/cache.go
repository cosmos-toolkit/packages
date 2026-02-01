// Package cache provides cache interface and in-memory implementation.
// Redis can be added as another implementation.
package cache

import (
	"context"
	"sync"
	"time"
)

// Cache interface Get/Set/Delete with TTL.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type item struct {
	value   []byte
	expires time.Time
}

// Memory implements Cache in memory with TTL.
type Memory struct {
	mu    sync.RWMutex
	items map[string]item
}

// NewMemory creates an in-memory cache.
func NewMemory() *Memory {
	c := &Memory{items: make(map[string]item)}
	go c.cleanup()
	return c
}

func (c *Memory) Get(ctx context.Context, key string) ([]byte, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || !it.expires.IsZero() && time.Now().After(it.expires) {
		if ok {
			c.Delete(ctx, key)
		}
		return nil, false
	}
	return it.value, true
}

func (c *Memory) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var expires time.Time
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}
	c.mu.Lock()
	c.items[key] = item{value: value, expires: expires}
	c.mu.Unlock()
	return nil
}

func (c *Memory) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *Memory) cleanup() {
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for range tick.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if !v.expires.IsZero() && now.After(v.expires) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
