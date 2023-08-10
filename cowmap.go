package eventbus

import (
	"sync"
	"sync/atomic"
)

// CowMap is a wrapper of Copy-On-Write map
//
// If a fully meaningful CowMap is implemented, both sync.Map and
// CowMap utilize atomic.Value atomic operations to access the map
// during data reading, resulting in similar read performance.
// In reality, sync.Map is already a read-write separated structure,
// yet it has better write performance. Therefore, CowMap directly
// utilizes sync.Map as its internal structure.
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
