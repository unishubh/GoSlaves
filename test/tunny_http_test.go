package main

import (
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/valyala/fasthttp"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkTunnyHTTP(b *testing.B) {
	ln, err := net.Listen("tcp4", ":6666")
	if err != nil {
		panic(err)
	}

	server := fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.Write([]byte("Served"))
		},
	}

	go func() {
		sp, _ := tunny.CreatePool(120, func(obj interface{}) interface{} {
			server.ServeConn(obj.(net.Conn))
			return nil
		}).Open()
		defer sp.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				time.Sleep(time.Second)
			} else {
				sp.SendWorkAsync(conn, nil)
			}
		}
	}()

	var ok, ot, er, timeout uint32
	var wg sync.WaitGroup
	now := time.Now()
	for i := 0; i < 20000; i++ {
		c := fasthttp.Client{}

		wg.Add(1)
		go func() {
			sc, _, err := c.GetTimeout(nil, "http://localhost:6666", time.Second*2)
			if err != nil {
				atomic.AddUint32(&timeout, 1)
			} else {
				switch {
				case sc < 300:
					atomic.AddUint32(&ok, 1)
				case sc > 499:
					atomic.AddUint32(&er, 1)
				default:
					atomic.AddUint32(&ot, 1)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("200: %d\n500: %d\nOther: %d\nTimed out: %d\n\nTime: %v",
		ok, er, ot, timeout, time.Since(now)))
	ln.Close()
}
