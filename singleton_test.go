package eventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SingletonSubscribe(t *testing.T) {
	ResetSingleton()
	err := Subscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)
	assert.NotNil(t, singleton)

	err = Subscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)
	err = Subscribe("testtopic", 1)
	assert.Equal(t, ErrHandlerIsNotFunc, err)

	err = Subscribe("testtopic", func(topic string) error {
		return nil
	})
	assert.Equal(t, ErrHandlerParamNum, err)
	err = Subscribe("testtopic", func(topic int, payload int) error {
		return nil
	})

	assert.Equal(t, ErrHandlerFirstParam, err)
	singleton.Close()
	err = Unsubscribe("testtopic", busHandlerTwo)
	assert.Equal(t, ErrChannelClosed, err)
	Close()
}

func Test_SingletonUnsubscribe(t *testing.T) {
	ResetSingleton()
	err := Unsubscribe("testtopic", busHandlerOne)
	assert.Equal(t, ErrNoSubscriber, err)
	assert.NotNil(t, singleton)

	err = Subscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)

	err = Unsubscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)
	singleton.Close()

	err = Unsubscribe("testtopic", busHandlerTwo)
	assert.Equal(t, ErrChannelClosed, err)
	Close()
}

func Test_SingletonPublish(t *testing.T) {
	ResetSingleton()
	err := Publish("testtopic", 1)
	assert.Nil(t, err)
	assert.NotNil(t, singleton)

	err = Subscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				err := Publish("testtopic", i)
				assert.Nil(t, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	Close()
}

func Test_SingletonPublishSync(t *testing.T) {
	ResetSingleton()
	err := Publish("testtopic", 1)
	assert.Nil(t, err)
	assert.NotNil(t, singleton)

	err = Subscribe("testtopic", busHandlerOne)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				err := PublishSync("testtopic", i)
				assert.Nil(t, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	Close()
}

func BenchmarkSingletonPublish(b *testing.B) {
	ResetSingleton()
	Subscribe("testtopic", busHandlerOne)

	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			Publish("testtopic", i)
		}
		wg.Done()
	}()
	wg.Wait()
	Close()
}

func BenchmarkSingletonPublishSync(b *testing.B) {
	ResetSingleton()
	Subscribe("testtopic", busHandlerOne)

	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			PublishSync("testtopic", i)
		}
		wg.Done()
	}()
	wg.Wait()
	Close()
}
