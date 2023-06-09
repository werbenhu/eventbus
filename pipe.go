package eventbus

import (
	"sync"
)

type Handler[T any] func(payload T)

// Pipe is a wrapper for a channel that allows for asynchronous message passing to subscribers.
// Use Pipe.Publish() instead of `chan<-` and Pipe.Subscribe() instead of `<-chan`.
// To pass messages to subscribers synchronously, use Pipe.PublishSync(), which does not use a channel.
// If multiple subscribers exist, each subscriber will receive the message.
type Pipe[T any] struct {
	sync.RWMutex
	bufferSize int
	channel    chan T
	handlers   *CowMap
	closed     bool
	stopCh     chan any
}

// NewPipe create a unbuffered pipe
func NewPipe[T any]() *Pipe[T] {
	p := &Pipe[T]{
		bufferSize: -1,
		channel:    make(chan T),
		stopCh:     make(chan any),
		handlers:   NewCowMap(),
	}

	go p.loop()
	return p
}

// NewPipe create a buffered pipe, bufferSize is the buffer size of the pipe
// When create a buffered pipe. You can publish into the Pipe without a corresponding concurrent subscriber.
func NewBufferedPipe[T any](bufferSize int) *Pipe[T] {
	if bufferSize <= 0 {
		bufferSize = 1
	}

	p := &Pipe[T]{
		bufferSize: bufferSize,
		channel:    make(chan T, bufferSize),
		stopCh:     make(chan any),
		handlers:   NewCowMap(),
	}

	go p.loop()
	return p
}

// loop loops forever, receiving published message from the pipe, transfer payload to subscriber by calling handlers
func (p *Pipe[T]) loop() {
	for {
		select {
		case payload := <-p.channel:
			p.handlers.Range(func(key any, fn any) bool {
				fn.(Handler[T])(payload)
				return true
			})
		case <-p.stopCh:
			return
		}
	}
}

// subscribe add a handler to a pipe, return error if the pipe is closed.
func (p *Pipe[T]) Subscribe(handler Handler[T]) error {
	p.RLock()
	defer p.RUnlock()
	if p.closed {
		return ErrChannelClosed
	}
	p.handlers.Store(&handler, handler)
	return nil
}

// unsubscribe removes handler defined for this pipe.
func (p *Pipe[T]) Unsubscribe(handler Handler[T]) error {
	p.RLock()
	defer p.RUnlock()
	if p.closed {
		return ErrChannelClosed
	}
	p.handlers.Delete(&handler)
	return nil
}

// Publish triggers the handlers defined for this pipe, transferring the payload to the handlers.
func (p *Pipe[T]) Publish(payload T) error {
	p.RLock()
	defer p.RUnlock()
	if p.closed {
		return ErrChannelClosed
	}
	p.channel <- payload
	return nil
}

// PublishSync triggers the handlers defined for this pipe synchronously, without using a channel.
// The payload will be passed directly to the handlers.
func (p *Pipe[T]) PublishSync(payload T) error {
	p.RLock()
	defer p.RUnlock()
	if p.closed {
		return ErrChannelClosed
	}
	p.handlers.Range(func(key any, fn any) bool {
		fn.(Handler[T])(payload)
		return true
	})
	return nil
}

// close closes the pipe
func (p *Pipe[T]) Close() {
	p.Lock()
	defer p.Unlock()

	if p.closed {
		return
	}
	p.closed = true
	p.stopCh <- struct{}{}
	close(p.channel)
}
