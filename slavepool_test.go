package slaves

import (
	"testing"
	"time"
)

func TestServe_SlavePool(t *testing.T) {
	ch := make(chan int, 20)
	cs := make(chan struct{})
	sp := &SlavePool{
		Work: func(obj interface{}) {
			ch <- obj.(int)
		},
	}
	sp.Open()
	defer sp.Close()

	go func() {
		p := 0
		for range ch {
			p++
		}
		if p == 20 {
			cs <- struct{}{}
		} else {
			t.Fatal("Bad test: ", p)
		}
	}()

	for i := 0; i < 20; i++ {
		sp.Serve(i)
	}
	time.Sleep(time.Second)
	close(ch)

	select {
	case <-cs:
	case <-time.After(time.Second * 2):
		t.Fatal("timeout")
	}
}

func BenchmarkSlavePool(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})

	sp := &SlavePool{
		Work: func(obj interface{}) {
			ch <- obj.(int)
		},
		Size: b.N,
	}
	sp.Open()
	defer sp.Close()

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		done <- struct{}{}
		close(ch)
		close(done)
	}()

	for i := 0; i < b.N; i++ {
		sp.Serve(i)
	}
	<-done
}
