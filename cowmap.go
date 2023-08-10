package eventbus

import (
	"sync"
	"sync/atomic"
)

// CowMap is a wrapper of Copy-On-Write map
type CowMap struct {
	sync.Map
}

// CowMap creates a new CowMap instance
func NewCowMap() *CowMap {
	return &CowMap{}
}

// Len returns the number of key-value pairs stored in the map
func (c *CowMap) Len() uint32 {
	var size uint32
	c.Range(func(k, v any) bool {
		atomic.AddUint32(&size, 1)
		return true
	})
	return size
}

// Clear Removes all key-value pairs from the map
func (c *CowMap) Clear() {
	c.Range(func(k, v any) bool {
		c.Delete(k)
		return true
	})
}
