package slaves

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestServe_Queue(t *testing.T) {
	queue := DoQueue(5, func(obj interface{}) {
		fmt.Println(obj)
		time.Sleep(time.Second * 5)
	})
	defer queue.Close()

	files, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		for i := range files {
			queue.Serve(files[i].Name())
		}
	}
	time.Sleep(time.Second)
}
