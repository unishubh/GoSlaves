package slaves

import (
	"runtime"
	"testing"

	"github.com/Jeffail/tunny"
	"github.com/ivpusic/grpool"
	"github.com/themester/GoSlaves"
)

func BenchmarkGrPool(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})

	pool := grpool.NewPool(50, 50)
	defer pool.Release()

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

	for i := 0; i < b.N; i++ {
		pool.JobQueue <- func() {
			ch <- i
		}
	}
	<-done
	close(ch)
	close(done)
}

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

func BenchmarkSlavePool(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})

	sp := &slaves.SlavePool{
		Work: func(obj interface{}) {
			ch <- obj.(int)
		},
	}
	sp.Open()
	defer sp.Close()

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

	for i := 0; i < b.N; i++ {
		sp.Serve(i)
	}
	<-done
	close(ch)
	close(done)
}
