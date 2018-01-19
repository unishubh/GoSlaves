package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestServe_Queue(t *testing.T) {
	queue := DoQueue(5, func(obj interface{}) {
		fmt.Println(obj)
	})
	defer queue.WaitClose()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		for i := range files {
			queue.Serve(files[i].Name())
		}
	}
}
