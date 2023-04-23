package eventbus

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func pipeSubOne(val int) {
	// fmt.Printf("pipeSubOne:%d\n", val)
}

func pipeSubTwo(val int) {
	// fmt.Printf("pipeSubTwo:%d\n", val)
}

func Test_NewPipe(t *testing.T) {
	p := NewPipe[int]()
	assert.NotNil(t, p)
	assert.NotNil(t, p.channel)
	assert.NotNil(t, p.stopCh)
	assert.NotNil(t, p.handlers)
	p.Close()
}

func Test_PipeSubscribe(t *testing.T) {
	p := NewPipe[int]()
	assert.NotNil(t, p)
	assert.NotNil(t, p.channel)

	err := p.Subscribe(pipeSubOne)
	assert.Nil(t, err)
	p.Close()
	err = p.Subscribe(pipeSubTwo)
	assert.Equal(t, ErrChannelClosed, err)
}

func Test_PipeUnsubscribe(t *testing.T) {
	p := NewPipe[int]()
	assert.NotNil(t, p)
	assert.NotNil(t, p.channel)

	err := p.Subscribe(pipeSubOne)
	assert.Nil(t, err)
	err = p.Unsubscribe(pipeSubOne)
	assert.Nil(t, err)

	err = p.Subscribe(pipeSubOne)
	assert.Nil(t, err)
	p.Close()
	err = p.Unsubscribe(pipeSubOne)
	assert.Equal(t, ErrChannelClosed, err)
}

func Test_PipePublish(t *testing.T) {
	p := NewPipe[int]()
	assert.NotNil(t, p)
	assert.NotNil(t, p.channel)

	err := p.Subscribe(pipeSubOne)
	time.Sleep(time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 10000; i++ {
			err := p.Publish(i)
			assert.Nil(t, err)
		}
		wg.Done()
	}()
	wg.Wait()

	p.Close()
	err = p.Publish(1)
	assert.Equal(t, ErrChannelClosed, err)
}
