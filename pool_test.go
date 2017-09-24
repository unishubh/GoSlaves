package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
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

func TestSendWorkTo(t *testing.T) {
	sp := MakePool(2)

	sp.Open(func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer sp.Close()

	sp.Slaves[0].Type = []byte("Borja")
	sp.Slaves[1].Type = []byte("Paquillo")

	sp.SendWorkTo("Borja", "Make me a cake plsssss")
	sp.SendWorkTo("Paquillo", "Execute python and kill my motherboard")
}

func TestDeleteSlave(t *testing.T) {
	sp := MakePool(20)

	fmt.Println("Slaves:", sp.GetSlaves())

	sp.Open(func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer fmt.Println("Slaves:", sp.GetSlaves())
	defer sp.Close()

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		sp.DeleteSlave()
	}
}

func TestAddSlave(t *testing.T) {
	sp := MakePool(1)

	fmt.Println("Slaves:", sp.GetSlaves())

	sp.Open(func(obj interface{}) interface{} {
		fmt.Println(obj)
		return nil
	}, nil)
	defer fmt.Println("Slaves:", sp.GetSlaves())
	defer sp.Close()

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		sp.AddSlave()
	}
}

func TestMakePool(t *testing.T) {
	sp := MakePool(10)
	sp.Open(func(obj interface{}) interface{} {
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
