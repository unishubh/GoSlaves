package slaves

import (
	"sync"
)

type Works struct {
	lock     sync.RWMutex
	elements []interface{}
}

func (w *Works) Get() interface{} {
	w.lock.Lock()
	defer w.lock.Unlock()

	n := len(w.elements) - 1
	if n < 0 {
		return nil
	}

	e := w.elements[n]
	if e != nil {
		w.elements = w.elements[:n]
	}

	return e
}

func (w *Works) Put(e interface{}) {
	if w.elements == nil {
		w.elements = make([]interface{}, 0)
	}
	w.lock.Lock()
	w.elements = append(w.elements, e)
	w.lock.Unlock()
}
