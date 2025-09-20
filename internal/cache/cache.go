// Ey, this package handles all the cache stuff
// It's re important for performance, ya know what I mean?
package cache

import (
	"sync"
	"time"
)

// Item holds the value and when it expires, re simple no?
type item struct {
	value     interface{}
	expiresAt time.Time
}

// Cache is our main struct che, it's like a map but with some extra magic
// We use mutex to avoid any quilombo with concurrent access, everything super zarpado
type Cache struct {
	mu       sync.RWMutex        // Mutex to avoid que se rompa todo with concurrent access
	items    map[string]item     // The actual storage, nothing fancy viste
	done     chan struct{}       // Channel to tell the cleanup goroutine "che, time to go home"
	stopOnce sync.Once          // Makes sure we don't close things twice, would be alta cagada
}

// New creates a fresh cache, everything ready to rock
// Also starts the cleanup goroutine in the background, re piola
func New() *Cache {
	cache := &Cache{
		items: make(map[string]item),
		done:  make(chan struct{}),
	}
	go cache.startCleanup() // Launch the cleanup goroutine, super important eh!
	return cache
}

// Set puts something in the cache for a while
// Like when you leave the mate somewhere and grab it later, ya know?
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get tries to find stuff in the cache
// If it's expired or not there, returns false, re simple boludo
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false // Nah, not here che
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key) // This one's past its prime, delete it
		return nil, false
	}

	return item.value, true // Found it! Everything copado
}

// Delete removes something from the cache
// Like when your code is a desastre and you need to start fresh
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Close tells the cleanup goroutine "che, time to go home"
// Super important to call this or you'll leave goroutines hanging like dirty ropa
func (c *Cache) Close() error {
	c.stopOnce.Do(func() {
		close(c.done) // Send the signal just once, no seas ansioso
	})
	return nil
}

// startCleanup runs in background, cleaning old stuff every 5 minutes
// It's like having someone pick up your empty mate cups while you code
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // Always cleanup after yourself, no seas croto

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.done:
			return // Time to go home che, cleanup is done
		}
	}
}

// cleanup removes all the expired items from the cache
// Like throwing out yesterday's pizza, ya know what I mean?
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiresAt) {
			delete(c.items, key) // This one's old, che. Get rid of it
		}
	}
}