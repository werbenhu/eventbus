package eventbus

import (
	"sync"
	"sync/atomic"
)

type tmap map[any]any

// CowMap is a wrapper of Copy-On-Write map
type CowMap struct {
	mu       sync.Mutex
	readable atomic.Value
}

func NewCowMap() *CowMap {
	m := make(tmap)
	c := &CowMap{}
	c.readable.Store(m)
	return c
}

// clone create a copy of the map
func (c *CowMap) clone() tmap {
	m := make(tmap)
	for k, v := range c.readable.Load().(tmap) {
		m[k] = v
	}
	return m
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (c *CowMap) Load(key any) (value any, ok bool) {
	value, ok = c.readable.Load().(tmap)[key]
	return
}

// Len returns how many values stored in the map.
func (c *CowMap) Len() int {
	return len(c.readable.Load().(tmap))
}

// Store sets the value for a key.
func (c *CowMap) Store(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()
	copy[key] = value
	c.readable.Store(copy)
}

// Delete deletes the value for a key.
func (c *CowMap) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()
	delete(copy, key)

	c.readable.Store(copy)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (c *CowMap) Range(f func(key, value any) bool) {
	for k, v := range c.readable.Load().(tmap) {
		if !f(k, v) {
			break
		}
	}
}
