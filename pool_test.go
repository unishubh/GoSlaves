package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSendWork_SlavePool(t *testing.T) {
	sp := &SlavePool{
		Work: func(obj interface{}) interface{} {
			fmt.Println(obj)
			return nil
		},
	}
	sp.Open()
	defer sp.Close()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		fmt.Println("Files in temp directory:")
		for i := range files {
			sp.SendWork(files[i].Name())
		}
	}
}
