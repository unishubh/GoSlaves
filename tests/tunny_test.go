package slaves

import (
	"runtime"
	"testing"

	"github.com/Jeffail/tunny"
)

func BenchmarkTunny(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})
	numCPUs := runtime.NumCPU()

	go func() {
		var i = 0
		for i < b.N {
			select {
			case <-ch:
				i++
			}
		}
		done <- struct{}{}
	}()

	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		ch <- payload.(int)
		return nil
	})
	defer pool.Close()

	for i := 0; i < b.N; i++ {
		pool.Process(i)
	}

	<-done
	close(ch)
	close(done)
}
