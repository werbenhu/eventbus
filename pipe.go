package eventbus

import (
	"sync"
)

type Handler[T any] func(payload T)

type Pipe[T any] struct {
	bufferSize int
	ch         chan T
	readers    sync.Map
}

func NewPipe[T any]() *Pipe[T] {
	p := &Pipe[T]{
		bufferSize: -1,
		ch:         make(chan T),
	}

	go p.loop()
	return p
}

func NewBufferedPipe[T any](bufferSize int) *Pipe[T] {
	if bufferSize <= 0 {
		bufferSize = 1
	}

	p := &Pipe[T]{
		bufferSize: bufferSize,
		ch:         make(chan T, bufferSize),
	}

	go p.loop()
	return p
}

func (p *Pipe[T]) loop() {
	for {
		payload := <-p.ch
		p.readers.Range(func(key any, fn any) bool {
			fn.(Handler[T])(payload)
			return true
		})
	}
}

func (p *Pipe[T]) Subscribe(reader Handler[T]) error {
	p.readers.Store(&reader, reader)
	return nil
}

func (p *Pipe[T]) Publish(payload T) {
	p.ch <- payload
}

func (p *Pipe[T]) Close() {
	if p.ch != nil {
		close(p.ch)
		p.ch = nil
	}
}
