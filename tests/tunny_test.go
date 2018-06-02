package slaves

import (
	"runtime"
	"testing"

	"github.com/Jeffail/tunny"
)

func BenchmarkTunny(b *testing.B) {
	ch := make(chan int, b.N)
	numCPUs := runtime.NumCPU() + 1

	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		ch <- payload.(int)
		return nil
	})

	go func() {
		for i := 0; i < b.N; i++ {
			pool.Process(i)
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
	pool.Close()
}
