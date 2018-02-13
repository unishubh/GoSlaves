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
		for i := 0; i < b.N; i++ {
			<-ch
		}
		done <- struct{}{}
		close(ch)
		close(done)
	}()

	for i := 0; i < b.N; i++ {
		pool.JobQueue <- func() {
			ch <- i
		}
	}
	<-done
}

func BenchmarkTunny(b *testing.B) {
	ch := make(chan int, b.N)
	done := make(chan struct{})
	numCPUs := runtime.NumCPU()

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		done <- struct{}{}
		close(ch)
		close(done)
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
