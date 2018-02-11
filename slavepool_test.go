package slaves

import (
	"net"
	"sync"
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

func BenchmarkServe_SlavePool(b *testing.B) {
	ln, err := net.Listen("tcp4", ":6666")
	if err != nil {
		panic(err)
	}

	go func() {
		sp := &SlavePool{
			Work: func(obj interface{}) {
				conn := obj.(net.Conn)
				data := make([]byte, 20)

				conn.Read(data)
				conn.Write([]byte("Hello world"))
				conn.Close()
			},
		}
		sp.Open()
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}

			sp.Serve(conn)
		}
	}()

	var wg sync.WaitGroup
	for t := 0; t < 5; t++ {
		for i := 0; i < 500; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn, err := net.Dial("tcp4", "127.0.0.1:6666")
				if err != nil {
					panic(err)
				}
				data := make([]byte, 20)

				conn.Write([]byte("Hello guys"))
				conn.Read(data)
				conn.Close()
			}()
		}
		wg.Wait()
	}
}
