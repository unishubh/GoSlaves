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
		time.Sleep(time.Second * 4)
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

func TestStop_Queue(t *testing.T) {
	queue := DoQueue(5, func(obj interface{}) {
		fmt.Println(obj)
	})
	defer queue.Close()

	for i := 0; i < 20; i++ {
		if i == 5 {
			queue.Stop()
			time.Sleep(time.Second)
		}
		queue.Serve(i)
	}
	time.Sleep(time.Second * 2)
	queue.Resume()
}
