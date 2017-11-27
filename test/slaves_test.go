package main

import (
	"github.com/themester/GoSlaves"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkSlavesReq(b *testing.B) {
	var requests uint32

	ln, err := net.Listen("tcp4", ":6666")
	if err != nil {
		panic(err)
	}

	atomic.StoreUint32(&requests, 0)

	go func() {
		sp := &slaves.SlavePool{
			Work: func(obj interface{}) {
				conn := obj.(net.Conn)
				defer conn.Close()

				conn.SetWriteDeadline(time.Now().Add(time.Second))

				if _, err := conn.Write([]byte("Hello guys")); err != nil {
					return
				}

				atomic.AddUint32(&requests, 1)
			},
		}
		sp.Open()
		defer sp.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			sp.Serve(conn)
		}
	}()

	var wg sync.WaitGroup
	for p := 0; p < 500; p++ {
		for i := 0; i < 400; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				conn, err := net.Dial("tcp4", "127.0.0.1:6666")
				if err != nil {
					return
				}
				defer conn.Close()

				conn.SetReadDeadline(time.Now().Add(time.Second))

				if _, err := ioutil.ReadAll(conn); err != nil {
					return
				}
			}()
		}
	}
	wg.Wait()

	println("served", atomic.LoadUint32(&requests))
}
