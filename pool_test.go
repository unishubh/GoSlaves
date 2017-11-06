package slaves

import (
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
)

func TestSendWork_SlavePool(t *testing.T) {
	sp := &SlavePool{
		Work: func(obj interface{}) {
			if obj == nil {
				panic("is nil")
			}
		},
	}
	sp.Open()
	defer sp.Close()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		for i := range files {
			sp.SendWork(files[i].Name())
		}
	}
}

func BenchmarkSendWork_SlavePool(b *testing.B) {
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
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}

			sp.SendWork(conn)
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
