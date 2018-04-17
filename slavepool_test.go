package slaves

import (
	"testing"
)

func TestServe_SlavePool(t *testing.T) {
	ch := make(chan int, 20)
	sp := NewPool(func(obj interface{}) {
		ch <- obj.(int)
	})

	rounds := 100000

	go func() {
		for i := 0; i < rounds; i++ {
			sp.Serve(i)
		}
	}()

	i := 0
	for i < rounds {
		select {
		case <-ch:
			i++
		}
	}
}

func BenchmarkSlavePool(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})

	sp := NewPool(func(obj interface{}) {
		ch <- obj.(int)
	})

	go func() {
		i := 0
		for i < b.N {
			select {
			case <-ch:
				i++
			}
		}
		done <- struct{}{}
	}()

	for i := 0; i < b.N; i++ {
		sp.Serve(i)
	}
	<-done
	close(ch)
	close(done)
}
