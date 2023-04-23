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

// NewCowMap creates a new CowMap instance
func NewCowMap() *CowMap {
	m := make(tmap)
	c := &CowMap{}
	c.readable.Store(m)
	return c
}

// clone creates a copy of the map by iterating over the original map
// and copying its key-value pairs to the new map
func (c *CowMap) clone() tmap {
	m := make(tmap)
	for k, v := range c.readable.Load().(tmap) {
		m[k] = v
	}
	return m
}

// Load returns the value stored in the map for a given key,
// or nil if the key is not present.
// The ok result indicates whether the value was found in the map.
func (c *CowMap) Load(key any) (value any, ok bool) {
	value, ok = c.readable.Load().(tmap)[key]
	return
}

// Len returns the number of key-value pairs stored in the map
func (c *CowMap) Len() int {
	return len(c.readable.Load().(tmap))
}

// Store sets the value for a given key by creating a new copy of the map
// and adding the new key-value pair to it
func (c *CowMap) Store(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()      // create a copy of the map
	copy[key] = value      // add the new key-value pair to the copy
	c.readable.Store(copy) // update the atomic value with the new copy
}

// Delete removes a key-value pair from the map by creating a new copy of the map
// and deleting the specified key from it
func (c *CowMap) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()
	delete(copy, key)

	c.readable.Store(copy)
}

// Range calls the provided function for each key-value pair in the map,
// stopping the iteration if the function returns false
func (c *CowMap) Range(f func(key, value any) bool) {
	for k, v := range c.readable.Load().(tmap) {
		if !f(k, v) {
			break
		}
	}
}
