package eventbus

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCowMap(t *testing.T) {
	m := NewCowMap()
	assert.NotNil(t, m)
}

func Test_CowMapLoad(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	for i := 0; i < 100; i++ {
		val, ok := m.Load(i)
		assert.True(t, ok)
		assert.Equal(t, strconv.Itoa(i), val.(string))
	}

	val, ok := m.Load(999999)
	assert.False(t, ok)
	assert.Equal(t, nil, val)
}

func Test_CowMapStore(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	assert.Equal(t, uint32(100), m.Len())
}

func Test_CowMapDelete(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	for i := 0; i < 50; i++ {
		m.Delete(i)
	}

	for i := 0; i < 50; i++ {
		val, ok := m.Load(i)
		assert.False(t, ok)
		assert.Equal(t, nil, val)
	}

	for i := 50; i < 100; i++ {
		val, ok := m.Load(i)
		assert.True(t, ok)
		assert.Equal(t, strconv.Itoa(i), val.(string))
	}
}

func Test_CowMapClear(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	m.Clear()
	assert.Equal(t, uint32(0), m.Len())
}

func Test_CowMapLen(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}
	assert.Equal(t, uint32(100), m.Len())
}

func Test_CowMapRange(t *testing.T) {
	m := NewCowMap()
	expert := make([]bool, 100)

	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}
	m.Range(func(key any, value any) bool {
		if value.(string) == strconv.Itoa(key.(int)) {
			expert[key.(int)] = true
		}
		return true
	})

	for _, val := range expert {
		assert.True(t, val)
	}
}

func Test_CowMapRangeStop(t *testing.T) {
	m := NewCowMap()
	results := make([]bool, 100)

	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	count := 0
	m.Range(func(key any, value any) bool {
		// only range ten elements and then stop range
		if count >= 10 {
			return false
		}

		count++
		if value.(string) == strconv.Itoa(key.(int)) {
			results[key.(int)] = true
		}
		return true
	})

	expert := 0
	for _, val := range results {
		if val {
			expert++
		}
	}
	assert.Equal(t, count, expert)
}

func Test_CowMapConcurrentLoadPanic(t *testing.T) {
	m := NewCowMap()
	assert.NotPanics(t, func() {
		for i := 0; i < 100; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					m.Store(j, strconv.Itoa(j))
				}
			}()
		}
	})
}

func Test_CowMapConcurrentStorePanic(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, strconv.Itoa(i))
	}

	assert.NotPanics(t, func() {
		for i := 0; i < 100; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					m.Load(j)
				}
			}()
		}
	})
}

func Test_CowMapStoreOrLoadConcurrent(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}

	storewg := sync.WaitGroup{}
	storewg.Add(100)
	assert.NotPanics(t, func() {
		for i := 0; i < 100; i++ {
			go func(index int) {
				for j := index * 100; j < (index+1)*100; j++ {
					m.Store(j, j)
				}
				storewg.Done()
			}(i)
		}
	})
	storewg.Wait()

	loadwg := sync.WaitGroup{}
	loadwg.Add(100)
	assert.NotPanics(t, func() {
		for i := 0; i < 100; i++ {
			go func(index int) {
				for j := index * 100; j < (index+1)*100; j++ {
					val, ok := m.Load(j)
					assert.True(t, ok)
					assert.Equal(t, j, val)
				}
				loadwg.Done()
			}(i)
		}
	})
	loadwg.Wait()
}

func Test_CowMapStoreAndLoadConcurrent(t *testing.T) {
	m := NewCowMap()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}

	assert.NotPanics(t, func() {
		loadGoroutineSize := 100
		loadWg := sync.WaitGroup{}
		loadWg.Add(loadGoroutineSize)

		for i := 0; i < loadGoroutineSize; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					val, ok := m.Load(j)
					assert.True(t, ok)
					assert.Equal(t, j, val)
				}
				loadWg.Done()
			}()
		}

		storeGoroutineSize := 100
		storeWg := sync.WaitGroup{}
		storeWg.Add(storeGoroutineSize)
		for i := 0; i < storeGoroutineSize; i++ {
			go func(index int) {
				for j := 0; j < 100; j++ {
					m.Store(j, j)
				}
				storeWg.Done()
			}(i)
		}

		storeWg.Wait()
		loadWg.Wait()
	})
}
