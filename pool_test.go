package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSendWork_SlavePool(t *testing.T) {
	sp := MakePool(5, func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer sp.Close()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		fmt.Println("Files in temp directory:")
		for i := range files {
			sp.SendWork(files[i].Name())
		}
	}
}

func TestSendWorkTo_SlavePool(t *testing.T) {
	sp := MakePool(4, func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer sp.Close()

	sp.Slaves[0].Name = "borja"

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		fmt.Println("Files in temp directory:")
		for i := range files {
			sp.SendWorkTo("borja", files[i].Name())
		}
	}
}

func TestAdd_SlavePool(t *testing.T) {
	sp := MakePool(1, nil, nil)
	defer sp.Close()

	for i := 0; i < 20; i++ {
		sp.Add(Slave{})
		fmt.Println(sp.Len())
	}
}

func TestDel_SlavePool(t *testing.T) {
	sp := MakePool(1, nil, nil)
	defer sp.Close()

	for i := 0; i < 20; i++ {
		s := Slave{}
		sp.Add(s)
		sp.Del()
		fmt.Println(sp.Len())
	}
}
