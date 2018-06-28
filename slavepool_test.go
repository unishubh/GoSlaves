package slaves

import (
	"sync"
	"sync/atomic"
	"testing"
)

func executeServe(sp *SlavePool, rounds int) {
	for i := 0; i < rounds; i++ {
		sp.Serve(i)
	}
}

func TestServe_SlavePool(t *testing.T) {
	ch := make(chan int, 1)
	counter := uint32(0)
	sp := NewPool(0, func(obj interface{}) {
		atomic.AddUint32(&counter, 1)
		ch <- obj.(int)
	})

	rounds := 100000

	go executeServe(&sp, rounds)

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

func TestServeMassiveServe_SlavePool(t *testing.T) {
	ch := make(chan struct{}, 1)

	locker := sync.Mutex{}
	counter := 0
	sp := NewPool(0, func(obj interface{}) {
		locker.Lock()
		counter++
		locker.Unlock()
		ch <- struct{}{}
	})

	rounds := 100000

	go executeServe(&sp, rounds/4)
	go executeServe(&sp, rounds/4)
	go executeServe(&sp, rounds/4)
	go executeServe(&sp, rounds/4)

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

	go executeServe(&sp, b.N)

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
