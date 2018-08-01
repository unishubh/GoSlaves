package slaves

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func executeServe(p *Pool, rounds int) {
	for i := 0; i < rounds; i++ {
		p.Serve(i)
	}
}

func TestServeNonStop_SlavePool(t *testing.T) {
	sp := NewPool(2, func(_ interface{}) {
		time.Sleep(time.Second)
	})
	if sp.ServeNonStop(nil) == false {
		t.Fatal("non-expected returned value")
	}
	if sp.ServeNonStop(nil) == false {
		t.Fatal("non-expected returned value")
	}
	if sp.ServeNonStop(nil) == true {
		t.Fatal("non-expected returned value")
	}
}

func TestServe_SlavePool(t *testing.T) {
	ch := make(chan int, 1)
	counter := uint32(0)
	locker := sync.Mutex{}
	sp := NewPool(0, func(obj interface{}) {
		locker.Lock()
		c := counter + 1
		if t := atomic.AddUint32(&counter, 1); t != c {
			panic(fmt.Sprintf("%d<>%d does not match", t, c))
		}
		locker.Unlock()
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

	counter := uint32(0)
	locker := sync.Mutex{}
	sp := NewPool(0, func(obj interface{}) {
		locker.Lock()
		c := counter + 1
		if t := atomic.AddUint32(&counter, 1); t != c {
			panic(fmt.Sprintf("%d<>%d does not match", t, c))
		}
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
