# GoSlaves

GoSlaves is a simple golang's library which can handle wide list of tasks asynchronously and safely.

[![GoDoc](https://godoc.org/github.com/themester/GoSlaves?status.svg)](https://godoc.org/github.com/themester/GoSlaves)
[![Go Report Card](https://goreportcard.com/badge/github.com/themester/goslaves)](https://goreportcard.com/report/github.com/themester/goslaves)

![alt text](https://raw.githubusercontent.com/themester/GoSlaves/master/logo.png)

Installation
------------

```
$ go get -u -v -x github.com/themester/GoSlaves
```

Benchmark
---------

After a lot of benchmarks and the following enhancings of the package I got this results:

```
$ GOMAXPROCS=4 go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
BenchmarkGrPool-4      	10000000	       738 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-4   	20000000	       353 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-4       	 2000000	      4190 ns/op	      32 B/op	       2 allocs/op
```

```
$ GOMAXPROCS=2 go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
BenchmarkGrPool-2      	20000000	       717 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-2   	100000000	       212 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-2       	 5000000	      3142 ns/op	      32 B/op	       2 allocs/op
```

Example
-------
```go
package main

import (
  "fmt"
  "net"

  "github.com/themester/GoSlaves"
)

func main() {
  pool := slaves.NewPool(func(obj interface{}) {
    conn := obj.(net.Conn)
    fmt.Fprintf(conn, "Welcome to GoSlaves!\n")
    conn.Close()
  })

  ln, err := net.Listen("tcp4", ":8080")
  if err != nil {
    panic(err)
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      break
    }
    pool.Serve(conn)
  }
}
```
