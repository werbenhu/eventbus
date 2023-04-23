<div align='center'>
<a href="https://github.com/werbenhu/eventbus/actions"><img src="https://github.com/werbenhu/eventbus/workflows/Go/badge.svg"></a>
<a href="https://coveralls.io/github/werbenhu/eventbus?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/eventbus/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/eventbus"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
</div>

# EventBus
EventBus 是一个轻量级的事件总线，可以简化 Go 协程之间的通信。


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
EventBus 是一个事件主题的容器，每个主题对应一个通道。`eventbus.Publish()` 方法将消息推送到通道，`eventbus.Subscribe(`) 方法中的处理程序将处理从通道出来的消息。

### EventBus 示例
```go
func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {
	bus := eventbus.New()

	// Subscribe() 订阅一个主题，如果handler不是函数则返回错误。
	// handler必须有两个参数：第一个参数必须是字符串类型，
	// handler的第二个参数类型必须与 `Publish()` 中的 payload 类型一致。
	bus.Subscribe("testtopic", handler)

	// Publish() 方法触发为主题定义的handler。`payload` 参数将传递给handler。
	// payload 的类型必须与 `Subscribe()` 中handler的第二个参数类型相对应。
	bus.Publish("testtopic", 100)

	// 订阅者异步接收消息。为了确保订阅者可以接收所有消息，在取消订阅之前需要有一个延迟。
	time.Sleep(time.Millisecond)
	bus.Unsubscribe("testtopic", handler)
	bus.Close()
}
```

## 使用Pipe代替Channel

Pipe 是通道的一个封装。订阅者异步接收消息。您可以使用 `Pipe.Publish()` 方法代替 `chan <-`，使用 `Pipe.Subscribe()` 方法代替 `<-chan`。如果有多个订阅者，则每个订阅者将接收到发布出来的每一条消息。

如果要使用带缓冲的通道，可以使用 `eventbus.NewPipe[T](bufferSize int)` 方法创建带缓冲的管道。

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
