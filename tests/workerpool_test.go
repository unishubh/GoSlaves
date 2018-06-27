package slaves

import (
	"testing"

	"github.com/gammazero/workerpool"
)

func BenchmarkWorkerpool(b *testing.B) {
	ch := make(chan int, b.N)

	wp := workerpool.New(4)

	go func() {
		for i := 0; i < b.N; i++ {
			wp.Submit(func() {
				ch <- i
			})
		}
	}()

	var i = 0
	for i < b.N {
		select {
		case <-ch:
			i++
		}
	}

	close(ch)
	wp.Stop()
}
