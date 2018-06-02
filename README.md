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
BenchmarkGrPool-4      	20000000	       717 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-4   	50000000	       367 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-4       	 3000000	      4131 ns/op	      32 B/op	       2 allocs/op
```

```
$ GOMAXPROCS=2 go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
BenchmarkGrPool-2      	20000000	       740 ns/op	      40 B/op	       1 allocs/op
BenchmarkSlavePool-2   	50000000	       353 ns/op	      16 B/op	       1 allocs/op
BenchmarkTunny-2       	 5000000	      3370 ns/op	      32 B/op	       2 allocs/op
```

Optimizations
-------------

You can optimize this package changing ChanSize variable and GOMAXPROCS env var. Here is a lot of benchmarks using different sizes:

- ChanSize 100 and GOMAXPROCS 4:
```
BenchmarkSlavePool-4   	50000000	       263 ns/op	      16 B/op	       1 allocs/op
```

- ChanSize 1000 and GOMAXPROCS 4:
```
BenchmarkSlavePool-4   	100000000	       190 ns/op	      16 B/op	       1 allocs/op
```

- ChanSize 10000 and GOMAXPROCS 4:
```
BenchmarkSlavePool-4   	100000000	       192 ns/op	      16 B/op	       1 allocs/op
```

Example
-------
```go
package main

import (
  "github.com/themester/GoSlaves"
)

func main() {
  ch := make(chan int, 20)
  pool := slaves.NewPool(func(o interface{}) {
    ch <- o.(int)
  })

  go func() {
    for i := 0; i < 100000; i++ {
      pool.Serve(i)
    }
  }()

  i := 0
  for i < 100000 {
    select {
    case <-ch:
      i++
    }
  }
  pool.Close()
}
```
