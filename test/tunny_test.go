package main

import (
	"github.com/Jeffail/tunny"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkTunny(b *testing.B) {
	var requests uint32

	ln, err := net.Listen("tcp", ":6664")
	if err != nil {
		panic(err)
	}

	atomic.StoreUint32(&requests, 0)

	go func() {
		sp, _ := tunny.CreatePool(4, func(obj interface{}) interface{} {
			conn := obj.(net.Conn)
			defer conn.Close()

			if _, err := conn.Write([]byte("Hello guys")); err != nil {
				panic(err)
			}

			atomic.AddUint32(&requests, 1)
			return nil
		}).Open()
		defer sp.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			sp.SendWorkAsync(conn, nil)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(3000)
	for i := 0; i < 3000; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp4", "127.0.0.1:6664")
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			conn.SetReadDeadline(time.Now().Add(time.Second))

			data := make([]byte, 10)
			if _, err = conn.Read(data); err != nil {
				return
			}
			data = nil
		}()
	}
	wg.Wait()

	println("served", atomic.LoadUint32(&requests))
}
