package utils

import (
	"sync"
	"time"
)

// CacheEntry represents a cached value with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// SimpleCache is a thread-safe in-memory cache with TTL support
type SimpleCache struct {
	mu    sync.RWMutex
	items map[string]CacheEntry
}

// NewSimpleCache creates a new cache instance
func NewSimpleCache() *SimpleCache {
	cache := &SimpleCache{
		items: make(map[string]CacheEntry),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Set stores a value in the cache with TTL
func (c *SimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from the cache
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

// Delete removes a value from the cache
func (c *SimpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all values from the cache
func (c *SimpleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheEntry)
}

// cleanup periodically removes expired entries
func (c *SimpleCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.items {
			if now.After(entry.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Size returns the number of items in the cache
func (c *SimpleCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
