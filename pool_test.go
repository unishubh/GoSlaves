package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
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

func TestSend_SlavePool(t *testing.T) {
	sp := MakePool(5, func(obj interface{}) interface{} {
		fmt.Println(obj)
		time.Sleep(time.Second)
		return nil
	}, nil)
	defer sp.Close()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		fmt.Println("Files in temp directory:")
		for i := range files {
			if sp.Working() == sp.Len() {
				sp.Send(files[i].Name())
			} else {
				sp.SendWork(files[i].Name())
			}
		}
	}
}

func TestAdd_SlavePool(t *testing.T) {
	sp := MakePool(0, nil, nil)
	defer sp.Close()

	for i := 0; i < 20; i++ {
		sp.Add(Slave{})
		fmt.Println(sp.Len())
	}
}

func TestDel_SlavePool(t *testing.T) {
	sp := MakePool(0, nil, nil)
	defer sp.Close()

	for i := 0; i < 20; i++ {
		s := Slave{}
		sp.Add(s)
		sp.Del()
		fmt.Println(sp.Len())
	}
}

func TestAddUSend_SlavePool(t *testing.T) {
	sp := MakePool(0, nil, nil)
	defer sp.Close()

	for i := 0; i < 20; i++ {
		s := Slave{
			Work: func(obj interface{}) interface{} {
				fmt.Println(obj, "Pool len:", sp.Len())
				return nil
			},
		}
		sp.AddUSend(s, i)
		sp.Del()
	}
}

func TestNonStacked_SlavePool(t *testing.T) {
	sp := MakePool(5, func(obj interface{}) interface{} {
		time.Sleep(3 * time.Second)
		return nil
	}, nil)
	sp.NotStack = true
	defer sp.Close()

	for i := 0; i < 5; i++ {
		sp.SendWork(nil)
	}
	err := sp.SendWork(nil)
	if err == nil {
		panic("stacked not works")
	}
	fmt.Println(err)
}

func TestForceClose_SlavePool(t *testing.T) {
	sp := MakePool(4, func(obj interface{}) interface{} {
		println(obj.(string))
		return nil
	}, nil)

	c := 0
	go func() {
		for _, str := range []string{
			"hello", "I", "am", "going",
			"to", "kill", "the", "pool",
		} {
			sp.SendWork(str)
			c++
		}
	}()
	sp.ForceClose()

	if c == 7 {
		panic("error")
	}
}
