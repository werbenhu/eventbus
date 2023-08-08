package eventbus

import (
	"testing"
)

func BenchmarkPipePublishSync(b *testing.B) {
	pipe := NewPipe[int]()

	pipe.Subscribe(pipeHandlerOne)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipe.PublishSync(i)
	}
}

func BenchmarkPipePublish(b *testing.B) {
	pipe := NewPipe[int]()

	pipe.Subscribe(pipeHandlerOne)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipe.Publish(i)
	}
}

func BenchmarkPipeGoChannel(b *testing.B) {
	ch := make(chan int)

	go func() {
		for {
			val := <-ch
			pipeHandlerOne(val)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
}
