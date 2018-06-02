package slaves

import (
	"testing"

	"github.com/themester/GoSlaves"
)

func BenchmarkSlavePool(b *testing.B) {
	ch := make(chan int, b.N)

	sp := slaves.NewPool(func(obj interface{}) {
		ch <- obj.(int)
	})

	go func() {
		for i := 0; i < b.N; i++ {
			sp.W <- i
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
	sp.Close()
}
