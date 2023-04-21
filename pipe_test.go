package eventbus

import (
	"fmt"
	"testing"
	"time"
)

func pipesub1(val int) {
	fmt.Printf("val:%d\n", val)
}
func pipesub2(val int) {
	fmt.Printf("va2:%d\n", val)
}
func pipesub3(val int) {
	fmt.Printf("va3:%d\n", val)
}

func Test_pipePublish(t *testing.T) {
	pipe := NewPipe[int]()
	pipe.Subscribe(pipesub1)
	pipe.Subscribe(pipesub2)
	pipe.Subscribe(pipesub3)
	go func() {
		for i := 0; i < 100000; i++ {
			pipe.Publish(i)
		}
	}()

	time.Sleep(1000 * time.Millisecond)

	// ch.close()
	// ch.publish(13)
}
