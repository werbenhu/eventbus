<div align='center'>
<a href="https://github.com/werbenhu/eventbus/actions"><img src="https://github.com/werbenhu/eventbus/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/eventbus"><img src="https://goreportcard.com/badge/github.com/werbenhu/eventbus"></a>
<a href="https://coveralls.io/github/werbenhu/eventbus?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/eventbus/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/eventbus"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/eventbus"><img src="https://pkg.go.dev/badge/github.com/werbenhu/eventbus.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)
# eventbus

A lightweight eventbus that simplifies communication between goroutines, it supports synchronous and asynchronous message publishing.

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

EventBus supports both synchronous and asynchronous message publication. it uses a `Copy-On-Write` map to manage handlers and topics, so it is not recommended for use in scenarios with a large number of frequent subscriptions and unsubscriptions.

#### Asynchronous Way

In EventBus, each topic corresponds to a channel. The `Publish()` method pushes the message to the channel, and the handler in the `Subscribe()` method handles the message that comes out of the channel.If you want to use a buffered EventBus, you can create a buffered EventBus with the `eventbus.NewBuffered(bufferSize int)` method, which will create a buffered channel for each topic.

#### Synchronous Way

In the synchronous way, EventBus does not use channels, but passes payloads to subscribers by calling the handler directly. To publish messages synchronously, use the `eventbus.PublishSync()` function.

### eventbus example
```go
package main

import (
	"fmt"
	"time"

	"github.com/werbenhu/eventbus"
)

func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {
	bus := eventbus.New()

	// Subscribe to a topic. Returns an error if the handler is not a function.
	// The handler function must have two parameters: the first parameter must be of type string,
	// and the second parameter's type must match the type of `payload` in the `Publish()` function.
	bus.Subscribe("testtopic", handler)

	// Publish a message asynchronously.
	// The `Publish()` function triggers the handler defined for the topic, and passes the `payload` as an argument.
	// The type of `payload` must match the type of the second parameter in the handler function defined in `Subscribe()`.
	bus.Publish("testtopic", 100)

	// Publish a message synchronously.
	bus.PublishSync("testtopic", 200)

	// Wait a bit to ensure that subscribers have received all asynchronous messages before unsubscribing.
	time.Sleep(time.Millisecond)
	bus.Unsubscribe("testtopic", handler)

	// Close the event bus.
	bus.Close()
}

```

### Using the global singleton object of EventBus
To make it more convenient to use EventBus, there is a global singleton object for EventBus. The internal channel of this singleton is unbuffered, and you can directly use `eventbus.Subscribe()`, `eventbus.Publish()`, and `eventbus.Unsubscribe()` to call the corresponding methods of the singleton object.

```go
package main

import (
	"fmt"
	"time"

	"github.com/werbenhu/eventbus"
)

func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {

	// eventbus.Subscribe() will call the global singleton's Subscribe() method
	eventbus.Subscribe("testtopic", handler)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Asynchronously publish messages
		for i := 0; i < 100; i++ {
			// eventbus.Publish() will call the global singleton's Publish() method
			eventbus.Publish("testtopic", i)
		}
		// Synchronously publish messages
		for i := 100; i < 200; i++ {
			eventbus.PublishSync("testtopic", i)
		}
		wg.Done()
	}()
	wg.Wait()

	time.Sleep(time.Millisecond)
	// eventbus.Unsubscribe() will call the global singleton's Unsubscribe() method
	eventbus.Unsubscribe("testtopic", handler)

	// eventbus.Close() will call the global singleton's Close() method
	eventbus.Close()
}
```

## Use Pipe instead of channel

Pipe is a wrapper for a channel without the concept of topics, with the generic parameter corresponding to the type of the channel. `eventbus.NewPipe[T]()` is equivalent to `make(chan T)`. Publishers publish messages, and subscribers receive messages. You can use the `Pipe.Publish()` method instead of `chan <-`, and the `Pipe.Subscribe()` method instead of `<-chan`. 

If there are multiple subscribers, each subscriber will receive every message that is published.If you want to use a buffered channel, you can use the `eventbus.NewBufferedPipe[T](bufferSize int)` method to create a buffered pipe.Pipe also supports synchronous and asynchronous message publishing. If you need to use the synchronous method, call `Pipe.PublishSync()`.

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
	go func() {
		for i := 0; i < 100; i++ {
			pipe.Publish(strconv.Itoa(i))
		}
		for i := 100; i < 200; i++ {
			pipe.PublishSync(strconv.Itoa(i))
		}
		wg.Done()
	}()
	wg.Wait()

	time.Sleep(time.Millisecond)
	pipe.Unsubscribe(handler1)
	pipe.Unsubscribe(handler2)
	pipe.Close()
}
```
