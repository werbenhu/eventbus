package eventbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sub1(topic string, val int) {
	// fmt.Printf("sub1 topic:%s, val:%d\n", topic, val)
}

func sub2(topic string, val int) {
	// fmt.Printf("sub2 topic:%s, val:%d\n", topic, val)
}

func Test_newChannel(t *testing.T) {
	ch := newChannel("test_topic", -1)
	assert.NotNil(t, ch)
	assert.NotNil(t, ch.channel)
	assert.Equal(t, "test_topic", ch.topic)
	assert.NotNil(t, ch.stopCh)
	assert.NotNil(t, ch.handlers)
	ch.close()
}

func Test_channelSubscribe(t *testing.T) {
	ch := newChannel("test_topic", -1)
	assert.NotNil(t, ch)
	assert.NotNil(t, ch.channel)
	assert.Equal(t, "test_topic", ch.topic)

	err := ch.subscribe(sub1)
	assert.Nil(t, err)
	ch.close()
	err = ch.subscribe(sub2)
	assert.Equal(t, ErrChannelClosed, err)
}

func Test_channelUnsubscribe(t *testing.T) {
	ch := newChannel("test_topic", -1)
	assert.NotNil(t, ch)
	assert.NotNil(t, ch.channel)
	assert.Equal(t, "test_topic", ch.topic)

	err := ch.subscribe(sub1)
	assert.Nil(t, err)
	err = ch.unsubscribe(sub1)
	assert.Nil(t, err)

	err = ch.subscribe(sub1)
	assert.Nil(t, err)
	ch.close()
	err = ch.subscribe(sub2)
	assert.Equal(t, ErrChannelClosed, err)
}

func Test_channelPublish(t *testing.T) {
	ch := newChannel("test_topic", -1)
	assert.NotNil(t, ch)
	assert.NotNil(t, ch.channel)
	assert.Equal(t, "test_topic", ch.topic)

	ch.subscribe(sub1)
	time.Sleep(time.Millisecond)

	go func() {
		for i := 0; i < 10000; i++ {
			err := ch.publish(i)
			assert.Nil(t, err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	ch.close()
	err := ch.publish(1)
	assert.Equal(t, ErrChannelClosed, err)
}
