package slaves

import (
	"fmt"
	"os"
	"io/ioutil"
	"testing"
)

func TestSendWork(t *testing.T) {
	sp := MakePool(2)

	sp.Open(func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer sp.Close()

	sp.SendWork("Make me a cake plsssss")
	sp.SendWork("Execute python and kill my motherboard")
}

func TestMakePool(t *testing.T) {
	sp := MakePool(10)
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

