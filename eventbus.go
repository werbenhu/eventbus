package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

// channel is the box for a topic and handlers. a topic corresponds to a channel
type channel struct {
	sync.RWMutex
	bufferSize int
	topic      string
	channel    chan any
	handlers   *CMap
	closed     bool
	stopCh     chan any
}

// newChannel create a channle for a topic
func newChannel(topic string, bufferSize int) *channel {
	var ch chan any
	if bufferSize <= 0 {
		ch = make(chan any)
	} else {
		ch = make(chan any, bufferSize)
	}
	c := &channel{
		topic:      topic,
		bufferSize: bufferSize,
		channel:    ch,
		handlers:   NewCMap(),
		stopCh:     make(chan any),
	}
	go c.loop()
	return c
}

// loop loops forever, receiving published message from the channel, transfer payload to subscriber by calling handlers
func (c *channel) loop() {
	topic := reflect.ValueOf(c.topic)
	for {
		select {
		case param := <-c.channel:
			c.handlers.Range(func(key any, fn any) bool {

				var payload reflect.Value
				handler := fn.(*reflect.Value)
				typ := handler.Type()

				if param == nil {
					payload = reflect.New(typ.In(1)).Elem()
				} else {
					payload = reflect.ValueOf(param)
				}

				payload = reflect.ValueOf(param)
				(*handler).Call([]reflect.Value{topic, payload})
				return true
			})
		case <-c.stopCh:
			return
		}
	}
}

// subscribe add a handler to a channel, return error if the channel is closed.
func (c *channel) subscribe(handler any) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return fmt.Errorf("channel on topic:%s is closed", c.topic)
	}
	fn := reflect.ValueOf(handler)
	c.handlers.Store(fn.Pointer(), &fn)
}

// publish trigger handlers defined for this channel. payload argument will be transferred to handlers.
func (c *channel) publish(payload any) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return fmt.Errorf("channel on topic:%s is closed", c.topic)
	}
	c.channel <- payload
	return nil
}

// unsubscribe removes handler defined for this channel.
func (c *channel) unsubscribe(handler any) {
	fn := reflect.ValueOf(handler)
	c.handlers.Delete(fn.Pointer())
}

// close closes a channel
func (c *channel) close() {
	c.Lock()
	defer c.Unlock()
	c.closed = true
	c.stopCh <- struct{}{}
	close(c.channel)
}

// EventBus is the box for topics and handlers.
type EventBus struct {
	channels   *CMap
	bufferSize int
	once       sync.Once
}

// NewBuffered returns new EventBus with a buffered channel.
// The second argument indicate the buffer's length
func NewBuffered(bufferSize int) *EventBus {
	if bufferSize <= 0 {
		bufferSize = 1
	}
	return &EventBus{
		bufferSize: bufferSize,
		channels:   NewCMap(),
	}
}

// New returns new EventBus with empty handlers.
func New() *EventBus {
	return &EventBus{
		bufferSize: -1,
	}
}

// Unsubscribe removes handler defined for a topic.
// Returns error if there are no handlers subscribed to the topic.
func (e *EventBus) Unsubscribe(topic string, handler any) error {
	ch, ok := e.channels.Load(topic)
	if !ok {
		return fmt.Errorf("no subscriber on topic:%s", topic)
	}
	ch.(*channel).unsubscribe(handler)
	return nil
}

// Subscribe subscribes to a topic, return error if handler is not a function.
func (e *EventBus) Subscribe(topic string, handler any) error {
	typ := reflect.TypeOf(handler)
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("the type of handler is %s, not type reflect.Func", reflect.TypeOf(handler).Kind())
	}
	if typ.NumIn() != 2 {
		return fmt.Errorf("the number of parameters of the handler must be two")
	}
	if typ.In(0).Kind() != reflect.String {
		return fmt.Errorf("the first of parameters of the handler must be string type")
	}

	ch, ok := e.channels.Load(topic)
	if !ok {
		ch = newChannel(topic, e.bufferSize)
		e.channels.Store(topic, ch)
	}
	ch.(*channel).subscribe(handler)
	return nil
}

// Publish trigger handlers defined for a topic. payload argument will be transferred to the handler.
func (e *EventBus) Publish(topic string, payload any) error {
	ch, ok := e.channels.Load(topic)

	if !ok {
		ch = newChannel(topic, e.bufferSize)
		e.channels.Store(topic, ch)
		go ch.(*channel).loop()
	}

	ch.(*channel).publish(payload)
	return nil
}

// Close closes the eventbus
func (e *EventBus) Close() {
	e.once.Do(func() {
		e.channels.Range(func(key any, ch any) bool {
			ch.(*channel).close()
			return true
		})
	})
}
