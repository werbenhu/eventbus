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

## Example

### eventbus example
```go
func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {
	bus := eventbus.New()
	bus.Subscribe("testtopic", handler)
	bus.Publish("testtopic", 100)

	time.Sleep(time.Millisecond)
	bus.Unsubscribe("testtopic", handler)
}
```

### pipe example
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

	time.Sleep(time.Millisecond)
	pipe.Unsubscribe(handler1)
	pipe.Unsubscribe(handler2)
}
```