package slaves

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestServe_SlavePool(t *testing.T) {
	ch := make(chan int, 1)
	counter := uint32(0)
	sp := NewPool(0, func(obj interface{}) {
		atomic.AddUint32(&counter, 1)
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
	if i != int(counter) {
		t.Fatalf("%d<>%d", i, rounds)
	}
	sp.Close()
}

func TestServeTimeout_SlavePool(t *testing.T) {
	ch := make(chan int, 1)

	counter := uint32(0)
	sp := NewPool(0, func(obj interface{}) {
		atomic.AddUint32(&counter, 1)
		time.Sleep(time.Second)
		ch <- obj.(int)
	})

	rounds := 10

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
	sp.Close()
	if i != int(counter) {
		t.Fatalf("%d<>%d", i, counter)
	}
}

func BenchmarkSlavePool(b *testing.B) {
	ch := make(chan int, b.N)

	sp := NewPool(0, func(obj interface{}) {
		ch <- obj.(int)
	})

	go func() {
		for i := 0; i < b.N; i++ {
			sp.Serve(i)
		}
	}()

	i := 0
	for i < b.N {
		select {
		case <-ch:
			i++
		}
	}
	close(ch)
	sp.Close()
}
