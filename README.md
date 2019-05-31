# GoSlaves

GoSlaves is a simple golang's library which can handle wide list of tasks asynchronously and safely.

[![GoDoc](https://godoc.org/github.com/dgrr/GoSlaves?status.svg)](https://godoc.org/github.com/dgrr/GoSlaves)
[![Go Report Card](https://goreportcard.com/badge/github.com/dgrr/goslaves)](https://goreportcard.com/report/github.com/dgrr/goslaves)

![alt text](https://raw.githubusercontent.com/dgrr/GoSlaves/master/logo.png)

Installation
------------

```
$ go get -u -v -x github.com/dgrr/GoSlaves
```

Benchmark
---------

Note that all of this benchmarks have been implemented as his owners recommends.
More of this goroutine pools works with more than 4 goroutines.

After a lot of benchmarks and the following enhancings of the package I got this results:

```
$ GOMAXPROCS=4 go test -v -bench=. -benchtime=5s -benchmem
goos: linux
goarch: amd64
BenchmarkGrPool-4       	10000000	       715 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-4    	20000000	       358 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-4        	 2000000	      4165 ns/op	      32 B/op	       2 allocs/op
BenchmarkWorkerpool-4   	 3000000	      3023 ns/op	      40 B/op	       1 allocs/op
```

```
$ GOMAXPROCS=2 go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
BenchmarkGrPool-2      	20000000	       717 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-2   	100000000	       212 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-2       	 5000000	      3142 ns/op	      32 B/op	       2 allocs/op
```

Library | Goroutines | Channel buffer
--- | --- | ---
GoSlaves | 4 | 1
GrPool | 50 | 50
Tunny | 4 | 1
Workerpool | 4 | 1

Example
-------
```go
package main

import (
  "fmt"
  "net"

  "github.com/dgrr/GoSlaves"
)

func main() {
  pool := slaves.NewPool(0, func(obj interface{}) {
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
