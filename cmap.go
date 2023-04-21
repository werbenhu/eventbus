package eventbus

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type innermap map[any]any

// CMap is the box of a Copy-On-Write map
type CMap struct {
	mu       sync.Mutex
	readable *innermap
}

func NewCMap() *CMap {
	inmap := make(innermap)
	return &CMap{
		readable: &inmap,
	}
}

// clone create a copy of the map
func (c *CMap) clone() innermap {
	m := make(innermap)
	for k, v := range *c.readable {
		m[k] = v
	}
	return m
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (c *CMap) Load(key any) (value any, ok bool) {
	value, ok = (*c.readable)[key]
	return
}

// Len returns how many values stored in the map.
func (c *CMap) Len() int {
	return len(*c.readable)
}

// Store sets the value for a key.
func (c *CMap) Store(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()
	copy[key] = value

	ptr := (*unsafe.Pointer)(unsafe.Pointer(&c.readable))
	atomic.SwapPointer(ptr, unsafe.Pointer(&copy))
}

// Delete deletes the value for a key.
func (c *CMap) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := c.clone()
	delete(copy, key)

	ptr := (*unsafe.Pointer)(unsafe.Pointer(&c.readable))
	atomic.SwapPointer(ptr, unsafe.Pointer(&copy))
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (c *CMap) Range(f func(key, value any) bool) {
	for k, v := range *c.readable {
		if !f(k, v) {
			break
		}
	}
}
