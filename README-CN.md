<div align='center'>
<a href="https://github.com/werbenhu/eventbus/actions"><img src="https://github.com/werbenhu/eventbus/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/eventbus"><img src="https://goreportcard.com/badge/github.com/werbenhu/eventbus"></a>
<a href="https://coveralls.io/github/werbenhu/eventbus?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/eventbus/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/eventbus"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/eventbus"><img src="https://pkg.go.dev/badge/github.com/werbenhu/eventbus.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)
# EventBus
EventBus 是一个轻量级的事件发布/订阅框架，支持同步和异步发布消息，它可以简化 Go 协程之间的通信。


## 安装

确保计算机上已安装 Go（版本 1.18+）。在终端中输入以下命令：

`go get github.com/werbenhu/eventbus`

*在项目中导入包*
```go
import (
	"github.com/werbenhu/eventbus"
)
```

## EventBus 是什么？

EventBus同时支持同步和异步的方式发布消息。EventBus使用一个Copy-On-Write的map管理handler和topic，所以不建议在有大量频繁的订阅和取消订阅的业务场景中使用。

#### 异步的方式
在EventBus里，每个主题对应一个通道。`Publish()` 方法将消息推送到通道，`Subscribe(`) 方法中的handler将处理从通道出来的消息。如果要使用带缓冲的EventBus，可以使用 `eventbus.NewBuffered(bufferSize int)` 方法创建带缓冲的EventBus，这样会为每个topic都创建一个带缓冲的channel。

#### 同步的方式
同步的方式下EventBus不使用channel，而是通过直接调用handler将消息传递给订阅者。如果想同步的方式发布消息，使用eventbus.PublishSync()函数即可。


### EventBus 示例
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

	// Subscribe() 订阅一个主题，如果handler不是函数则返回错误。
	// handler必须有两个参数：第一个参数必须是字符串类型，
	// handler的第二个参数类型必须与 `Publish()` 中的 payload 类型一致。
	bus.Subscribe("testtopic", handler)

	// 异步方式发布消息
	bus.Publish("testtopic", 100)

	// 同步方式发布消息
	bus.PublishSync("testtopic", 200)

	// 订阅者接收消息。为了确保订阅者可以接收完所有消息的异步消息，这里在取消订阅之前给了一点延迟。
	time.Sleep(time.Millisecond)
	bus.Unsubscribe("testtopic", handler)
	bus.Close()
}
```

### 使用全局的EventBus单例对象

为了更方便的使用EventBus, 这里有一个全局的EventBus单例对象，这个单例内部的channel是无缓冲的，直接使用`eventbus.Subscribe()`,`eventbus.Publish()`,`eventbus.Unsubscribe()`，将会调用该单例对象对应的方法。

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

	// eventbus.Subscribe() 将调用全局单例singleton.Subscribe()方法
	eventbus.Subscribe("testtopic", handler)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 异步方式发布消息
		for i := 0; i < 100; i++ {
			// 调用全局单例singleton.Publish()方法
			eventbus.Publish("testtopic", i)
		}
		// 同步方式发布消息
		for i := 100; i < 200; i++ {
			// 调用全局单例singleton.PublishSync()方法
			eventbus.PublishSync("testtopic", i)
		}
		wg.Done()
	}()
	wg.Wait()

	time.Sleep(time.Millisecond)
	// eventbus.Unsubscribe() 将调用全局单例singleton.Unsubscribe()方法
	eventbus.Unsubscribe("testtopic", handler)

	// eventbus.Close() 将调用全局单例singleton.Close()方法
	eventbus.Close()
}
```

## 使用Pipe代替Channel

Pipe 将通道封装成泛型对象，泛型参数对应channle里的类型，这里没有主题的概念。
`eventbus.NewPipe[T]()` 等价于 `make(chan T)`,发布者发布消息，订阅者接收消息，可以使用 `Pipe.Publish()` 方法代替 `chan <-`，使用 `Pipe.Subscribe()` 方法代替 `<-chan`。如果有多个订阅者，则每个订阅者将接收到发布出来的每一条消息。

如果要使用带缓冲的通道，可以使用 `eventbus.NewBufferedPipe[T](bufferSize int)` 方法创建带缓冲的管道。Pipe同样支持同步和异步的方式发布消息。如果需要使用同步的方式，请调用Pipe.PublishSync()。

#### Pipe 示例
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
