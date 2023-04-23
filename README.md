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
```
func handler(topic string, payload int) {
	fmt.Printf("topic:%s, payload:%d\n", topic, payload)
}

func main() {
	eb := eventbus.New()
	eb.Subscribe("testtopic", handler)
	eb.Publish("testtopic", 100)
	eb.Unsubscribe("testtopic", handler)
}
```

### pipe example
```
func handler1(payload string) {
	fmt.Printf("handler1 payload:%s\n", payload)
}

func handler2(payload string) {
	fmt.Printf("handler2 payload:%s\n", payload)
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
	pipe.Unsubscribe(handler1)
	pipe.Unsubscribe(handler2)
}
```