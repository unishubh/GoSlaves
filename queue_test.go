package slaves

import (
	"testing"
	"time"
)

func TestServe_Queue(t *testing.T) {
	ch := make(chan int, 10)
	done := make(chan struct{})
	queue := DoQueue(5, func(obj interface{}) {
		ch <- obj.(int)
		time.Sleep(time.Second * 1)
	})
	defer queue.Close()

	go func() {
		for t := 0; t < 10; t++ {
			<-ch
		}
		done <- struct{}{}
	}()

	for i := 0; i < 10; i++ {
		queue.Serve(i)
	}

	select {
	case <-done:
		close(ch)
		close(done)
		return
	case <-time.After(time.Second * 3):
		t.Fatal("timeout")
	}
}
