package slaves

import (
	"testing"

	"github.com/ivpusic/grpool"
)

func BenchmarkGrPool(b *testing.B) {
	ch := make(chan int, b.N)

	pool := grpool.NewPool(50, 50)
	defer pool.Release()

	go func() {
		for i := 0; i < b.N; i++ {
			pool.JobQueue <- func() {
				ch <- i
			}
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
}
