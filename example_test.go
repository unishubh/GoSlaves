package slaves_test

import (
	"fmt"
	"github.com/themester/GoSlaves"
	"io/ioutil"
	"os"
)

func ExampleSlavePool_SendWorkTo() {
	sp := slaves.MakePool(2)

	sp.Open(func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer sp.Close()

	sp.Slaves[0].Type = []byte("ProcessPiDecimals")
	sp.Slaves[1].Type = []byte("MakeCake")

	sp.SendWorkTo("MakeCake", "Make me a cake plsssss")
	sp.SendWorkTo("ProcessPiDecimals", "Execute python and kill my motherboard")
}

// Simple slave pool example
func ExampleMakePool() {
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
}
