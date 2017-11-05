# GoSlaves

GoSlaves is a simple golang's library which can handle wide list of tasks asynchronously

[![GoDoc](https://godoc.org/github.com/themester/GoSlaves?status.svg)](https://godoc.org/github.com/themester/GoSlaves)
[![Go Report Card](https://goreportcard.com/badge/github.com/themester/goslaves)](https://goreportcard.com/report/github.com/themester/goslaves)

![alt text](https://raw.githubusercontent.com/themester/GoSlaves/master/logo.png)

Installation
------------

Experimental version:
```bash
go get github.com/themester/GoSlaves
```

Stable version:
```bash
go get gopkg.in/themester/GoSlaves.v2
```

First stable version:
```bash
go get gopkg.in/themester/GoSlaves.v1
```

Example
-------
```go
package main

import (
  "fmt"
  "os"
  "time"
  "io/ioutil"
  "github.com/themester/GoSlaves"
)

func main() {
  sp := &slaves.SlavePool{
    Work: func(obj interface{}) interface{} {
      fmt.Println(obj)
      return nil
    },
  }
  sp.Open()
  defer func() {
    time.Sleep(time.Second)
    sp.Close()
  }()

  files, err := ioutil.ReadDir(os.TempDir())
  if err == nil {
    fmt.Println("Files in temp directory:")
    for i := range files {
      sp.SendWork(files[i].Name())
    }
  }
}
```
