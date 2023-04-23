package eventbus

import (
	"reflect"
	"sync"
)

// channel is the box of a topic and handlers. a topic corresponds to a channel
type channel struct {
	sync.RWMutex
	bufferSize int
	topic      string
	channel    chan any
	handlers   *CowMap
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
		handlers:   NewCowMap(),
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
		return ErrChannelClosed
	}
	fn := reflect.ValueOf(handler)
	c.handlers.Store(fn.Pointer(), &fn)
	return nil
}

// publish triggers the handlers defined for this channel. The `payload` argument will be passed to the handler.
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
	c.closed = true
	c.stopCh <- struct{}{}
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
	ch.(*channel).unsubscribe(handler)
	return nil
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
	ch.(*channel).subscribe(handler)
	return nil
}

// Publish triggers the handlers defined for a topic. The `payload` argument will be passed to the handler.
// The type of the payload must correspond to the second parameter of the handler in `Subscribe()`.
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
