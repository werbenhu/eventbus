package eventbus

import (
	"reflect"
	"sync"
)

// channel is a struct representing a topic and its associated handlers.
type channel struct {
	sync.RWMutex
	bufferSize int
	topic      string
	channel    chan any
	handlers   *CowMap
	closed     bool
	stopCh     chan any
}

// newChannel creates a new channel with a specified topic and buffer size.
// It initializes the handlers map with NewCowMap function and
// starts a goroutine c.loop() to continuously listen to messages in the channel.
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
		handlers:   NewCowMap(),
		stopCh:     make(chan any),
	}
	go c.loop()
	return c
}

// transfer calls all the handlers in the channel with the given payload.
// It iterates over the handlers in the handlers map to call them with the payload.
func (c *channel) transfer(topic string, payload any) {
	var payloadValue reflect.Value
	topicValue := reflect.ValueOf(c.topic)

	c.handlers.Range(func(key any, fn any) bool {
		handler := fn.(*reflect.Value)
		typ := handler.Type()

		if payload == nil {
			// If the parameter passed to the handler is nil,
			// it initializes a new payload element based on the
			// type of the second parameter of the handler using the reflect package.
			payloadValue = reflect.New(typ.In(1)).Elem()
		} else {
			payloadValue = reflect.ValueOf(payload)
		}
		(*handler).Call([]reflect.Value{topicValue, payloadValue})
		return true
	})
}

// loop listens to the channel and calls handlers with payload.
// It receives messages from the channel and then iterates over the handlers
// in the handlers map to call them with the payload.
func (c *channel) loop() {
	for {
		select {
		case payload := <-c.channel:
			c.transfer(c.topic, payload)
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
		return ErrChannelClosed
	}
	fn := reflect.ValueOf(handler)
	c.handlers.Store(fn.Pointer(), &fn)
	return nil
}

// publishSync triggers the handlers defined for this channel synchronously.
// The payload argument will be passed to the handler.
// It does not use channels and instead directly calls the handler function.
func (c *channel) publishSync(payload any) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return ErrChannelClosed
	}
	c.transfer(c.topic, payload)
	return nil
}

// publish triggers the handlers defined for this channel asynchronously.
// The `payload` argument will be passed to the handler.
// It uses the channel to asynchronously call the handler.
func (c *channel) publish(payload any) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return ErrChannelClosed
	}
	c.channel <- payload
	return nil
}

// unsubscribe removes handler defined for this channel.
func (c *channel) unsubscribe(handler any) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return ErrChannelClosed
	}
	fn := reflect.ValueOf(handler)
	c.handlers.Delete(fn.Pointer())
	return nil
}

// close closes a channel
func (c *channel) close() {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.stopCh <- struct{}{}
	c.handlers.Clear()
	close(c.channel)
}

// EventBus is a container for event topics.
// Each topic corresponds to a channel. `eventbus.Publish()` pushes a message to the channel,
// and the handler in `eventbus.Subscribe()` will process the message coming out of the channel.
type EventBus struct {
	channels   *CowMap
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
		channels:   NewCowMap(),
	}
}

// New returns new EventBus with empty handlers.
func New() *EventBus {
	return &EventBus{
		bufferSize: -1,
		channels:   NewCowMap(),
	}
}

// Unsubscribe removes handler defined for a topic.
// Returns error if there are no handlers subscribed to the topic.
func (e *EventBus) Unsubscribe(topic string, handler any) error {
	ch, ok := e.channels.Load(topic)
	if !ok {
		return ErrNoSubscriber
	}
	return ch.(*channel).unsubscribe(handler)
}

// Subscribe subscribes to a topic, return an error if the handler is not a function.
// The handler must have two parameters: the first parameter must be a string,
// and the type of the handler's second parameter must be consistent with the type of the payload in `Publish()`
func (e *EventBus) Subscribe(topic string, handler any) error {
	typ := reflect.TypeOf(handler)
	if typ.Kind() != reflect.Func {
		return ErrHandlerIsNotFunc
	}
	if typ.NumIn() != 2 {
		return ErrHandlerParamNum
	}
	if typ.In(0).Kind() != reflect.String {
		return ErrHandlerFirstParam
	}

	ch, ok := e.channels.Load(topic)
	if !ok {
		ch = newChannel(topic, e.bufferSize)
		e.channels.Store(topic, ch)
	}
	return ch.(*channel).subscribe(handler)
}

// publish triggers the handlers defined for this channel asynchronously.
// The `payload` argument will be passed to the handler.
// It uses the channel to asynchronously call the handler.
// The type of the payload must correspond to the second parameter of the handler in `Subscribe()`.
func (e *EventBus) Publish(topic string, payload any) error {
	ch, ok := e.channels.Load(topic)

	if !ok {
		ch = newChannel(topic, e.bufferSize)
		e.channels.Store(topic, ch)
		go ch.(*channel).loop()
	}

	return ch.(*channel).publish(payload)
}

// publishSync triggers the handlers defined for this channel synchronously.
// The payload argument will be passed to the handler.
// It does not use channels and instead directly calls the handler function.
func (e *EventBus) PublishSync(topic string, payload any) error {
	ch, ok := e.channels.Load(topic)

	if !ok {
		ch = newChannel(topic, e.bufferSize)
		e.channels.Store(topic, ch)
		go ch.(*channel).loop()
	}

	return ch.(*channel).publishSync(payload)
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
