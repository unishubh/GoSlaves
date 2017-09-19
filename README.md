# GoSlaves

GoSlaves is a simple golang's library which can handle wide list of tasks asynchroniously

[![GoDoc](https://godoc.org/github.com/themester/GoSlaves?status.svg)](https://godoc.org/github.com/themester/GoSlaves)
[![Go Report Card](https://goreportcard.com/badge/github.com/themester/goslaves)](https://goreportcard.com/report/github.com/themester/goslaves)

Installation
------------

```bash
go get github.com/themester/GoSlaves
```

Example
-------
```go
  sp := slaves.MakePool(10)
  if err := sp.Open(func(obj interface{}) interface{} {
    fmt.Println(obj.(string))
    return nil
  }, nil); err != nil {
    panic(err)
  }
  defer sp.Close()

  files, err := ioutil.ReadDir(os.TempDir())
  if err == nil {
    fmt.Println("Files in temp directory:")
    for i := range files {
      sp.SendWork(files[i].Name())
    }
  }
```
