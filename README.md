<div align='center'>
<a href="https://github.com/werbenhu/eventbus/actions"><img src="https://github.com/werbenhu/eventbus/workflows/Go/badge.svg"></a>
<a href="https://coveralls.io/github/werbenhu/eventbus?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/eventbus/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/eventbus"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
</div>

English | [简体中文](README-CN.md)
# eventbus
A lightweight eventbus that simplifies communication between goroutines.


## Installation

Make sure that go(version 1.18+) is installed on your computer. 
Type the following command:

`go get github.com/werbenhu/eventbus`

*Import package in your project*
```go
import (
	"github.com/werbenhu/eventbus"
)
```

## What's eventbus?
 EventBus is a container for event topics. Each topic corresponds to a channel. `eventbus.Publish()` pushes a message to the channel, and the handler in `eventbus.Subscribe()` will process the message coming out of the channel.

### eventbus example
```go
func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {
	bus := eventbus.New()

	// Subscribe() subscribes to a topic, return an error if the handler is not a function.
	// The handler must have two parameters: the first parameter must be a string,
	// and the type of the handler's second parameter must be consistent with the type of the payload in `Publish()`
	bus.Subscribe("testtopic", handler)

	// Publish() triggers the handlers defined for a topic. The `payload` argument will be passed to the handler.
	// The type of the payload must correspond to the second parameter of the handler in `Subscribe()`.
	bus.Publish("testtopic", 100)

	// Subscribers receive messages asynchronously. 
	// To ensure that subscribers can receive all messages, there is a delay before unsubscribe
	time.Sleep(time.Millisecond)
	bus.Unsubscribe("testtopic", handler)
	bus.Close()
}
```

## Use Pipes instead of channels

Pipe is a wrapper for a channel. Subscribers receive messages asynchronously. You can use `Pipe.Publish()` instead of `chan <-` and `Pipe.Subscribe()` instead of `<- chan`. If there are multiple subscribers, one message will be received by each subscriber.

If you want to use a buffered channel, you can use `eventbus.NewPipe[T](bufferSize int)` to create a buffered pipe.

#### pipe example
```go
func handler1(val string) {
	fmt.Printf("handler1 val:%s\n", val)
}

func handler2(val string) {
	fmt.Printf("handler2 val:%s\n", val)
}

func main() {
	pipe := eventbus.NewPipe[string]()
	pipe.Subscribe(handler1)
	pipe.Subscribe(handler2)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(p *eventbus.Pipe[string]) {
		for i := 0; i < 100; i++ {
			p.Publish(strconv.Itoa(i))
		}
		wg.Done()
	}(pipe)
	wg.Wait()

	// Subscribers receive messages asynchronously. 
	// To ensure that subscribers can receive all messages, there is a delay before unsubscribe
	time.Sleep(time.Millisecond)
	pipe.Unsubscribe(handler1)
	pipe.Unsubscribe(handler2)
	pipe.Close()
}

```
