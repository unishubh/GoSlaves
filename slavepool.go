package slaves

import (
	"runtime"
)

// This library uses a queue system. See Serve function.
type slave struct {
	ch chan interface{}
}

func newSlave(w func(interface{})) (s slave) {
	s.ch = make(chan interface{}, 1)
	go func() {
		var job interface{}
		for job = range s.ch {
			w(job)
		}
	}()
	return s
}

func (s *slave) close() {
	close(s.ch)
}

// SlavePool
type SlavePool struct {
	sv []slave
	n  int
}

// NewPool creates SlavePool.
//
// if workers is 0 default workers will be created
// use workers var if you know what you are doing
func NewPool(workers int, w func(interface{})) (sp SlavePool) {
	if w == nil {
		return
	}
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}

	sp.n = workers
	sp.sv = make([]slave, sp.n, sp.n)
	for i := 0; i < sp.n; i++ {
		sp.sv[i] = newSlave(w)
	}
	return
}

// Serve sends work to goroutine pool
//
// If all slaves are busy this function will stop until any of this ends a task.
func (sp *SlavePool) Serve(w interface{}) {
	i := 0
	for {
		select {
		case sp.sv[i].ch <- w:
			return
		default: // channel is busy
			i++
			if i == sp.n {
				i = 0
			}
		}
	}
}

// ServeNonStop returns true if work have been sended to the goroutine pool.
//
// This function returns a state and does not block the workflow.
func (sp *SlavePool) ServeNonStop(w interface{}) bool {
	i := 0
	for i < sp.n {
		select {
		case sp.sv[i].ch <- w:
			return true
		default:
			i++
		}
	}
	return false
}

// Close closes the SlavePool
func (sp *SlavePool) Close() {
	for i := 0; i < sp.n; i++ {
		sp.sv[i].close()
	}
}
